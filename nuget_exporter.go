package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// 預設的目標框架優先順序
var frameworkPriority = []string{
	"netstandard2.0",
	"net45",
	"net46",
	"net47",
	"net48",
	"netstandard2.1",
}

func getUserInput(prompt string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func createPackageJson(packageName, version, outputPath string) error {
	normalizedName := "com.nuget." + strings.ToLower(strings.ReplaceAll(packageName, ".", "-"))
	packageJsonContent := fmt.Sprintf(`{
  "name": "%s",
  "displayName": "%s",
  "version": "%s",
  "unity": "2019.1",
  "description": "Auto-generated package for %s",
  "dependencies": {}
}`, normalizedName, packageName, version, packageName)

	packageJsonPath := filepath.Join(outputPath, "package.json")
	return os.WriteFile(packageJsonPath, []byte(packageJsonContent), 0644)
}

func createAsmdef(asmName, dllName, outputPath string) error {
	asmdefName := strings.ToLower(strings.ReplaceAll(asmName, ".", "-")) + "-asmdef"
	asmdefContent := fmt.Sprintf(`{
  "name": "%s",
  "overrideReferences": true,
  "precompiledReferences": [
    "%s"
  ],
  "autoReferenced": true,
  "noEngineReferences": false
}`, asmdefName, dllName)

	asmdefPath := filepath.Join(outputPath, asmdefName+".asmdef")
	return os.WriteFile(asmdefPath, []byte(asmdefContent), 0644)
}

// 從找到的frameworkDirs中，根據預設的優先順序自動選擇最合適的
func chooseFrameworkAuto(frameworkDirs []string) string {
	fwSet := make(map[string]bool)
	for _, fw := range frameworkDirs {
		fwSet[fw] = true
	}
	for _, preferred := range frameworkPriority {
		if fwSet[preferred] {
			return preferred
		}
	}
	sort.Strings(frameworkDirs)
	return frameworkDirs[0]
}

// 簡易產生一個假GUID (16 hex chars)
func genGUID() string {
	const hexChars = "0123456789abcdef"
	b := make([]byte, 32)
	for i := 0; i < 32; i++ {
		b[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(b)
}

// 為一個資產產生簡易meta檔案內容
// 注意：實務上可能需要依資產類型產生更合適的 meta 檔
func generateMeta(guid string) []byte {
	return []byte(fmt.Sprintf(`fileFormatVersion: 2
guid: %s
timeCreated: %d
licenseType: Free
DefaultImporter:
  externalObjects: {}
  userData: 
  assetBundleName: 
  assetBundleVariant: 
`, guid, time.Now().Unix()))
}

// 掃描 export/<packageName> 下所有檔案，將它們全部打包到 .unitypackage
// 檔案路徑轉換規則：
// export/<packageName>/Runtime/... => Assets/<packageName>/Runtime/...
// 同理如有其他子檔案，皆以此類推
// 每個檔案以 GUID/asset, GUID/asset.meta, GUID/pathname 表示
func createUnityPackageFromExport(exportDir, packageName string, outPackageName string) error {
	// 收集所有需要打包的檔案
	var files []string
	err := filepath.Walk(exportDir, func(path string, info os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 開啟 unitypackage 檔案
	outFile, err := os.Create(outPackageName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, f := range files {
		rel, err := filepath.Rel(exportDir, f)
		if err != nil {
			return err
		}
		// 將路徑轉換為 Assets/<packageName>/...
		unityPath := filepath.Join("Assets", packageName, rel)
		unityPath = filepath.ToSlash(unityPath)

		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		guid := genGUID()

		// 寫入 asset
		assetHeader := &tar.Header{
			Name: guid + "/asset",
			Mode: 0600,
			Size: int64(len(content)),
		}
		if err := tarWriter.WriteHeader(assetHeader); err != nil {
			return err
		}
		if _, err := tarWriter.Write(content); err != nil {
			return err
		}

		// 寫入 asset.meta
		metaContent := generateMeta(guid)
		metaHeader := &tar.Header{
			Name: guid + "/asset.meta",
			Mode: 0600,
			Size: int64(len(metaContent)),
		}
		if err := tarWriter.WriteHeader(metaHeader); err != nil {
			return err
		}
		if _, err := tarWriter.Write(metaContent); err != nil {
			return err
		}

		// 寫入 pathname
		pathnameBytes := []byte(unityPath)
		pathnameHeader := &tar.Header{
			Name: guid + "/pathname",
			Mode: 0600,
			Size: int64(len(pathnameBytes)),
		}
		if err := tarWriter.WriteHeader(pathnameHeader); err != nil {
			return err
		}
		if _, err := tarWriter.Write(pathnameBytes); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Welcome to the Interactive NuGet to Unity Package Exporter!")
	nugetPackageName := getUserInput("Enter the NuGet package name (e.g. Newtonsoft.Json)", "")
	packageVersion := getUserInput("Enter the package version (or leave empty for latest)", "")

	// 建立暫存目錄
	tempDir, err := os.MkdirTemp("", "nuget_temp")
	if err != nil {
		fmt.Printf("Failed to create temporary directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	exportPath := "./export"
	pluginPath := filepath.Join(exportPath, nugetPackageName)

	// nuget install
	args := []string{"install", nugetPackageName, "-OutputDirectory", tempDir}
	if packageVersion != "" && packageVersion != "latest" {
		args = append(args, "-Version", packageVersion)
	}

	fmt.Printf("\nDownloading %s (%s) to temporary directory...\n", nugetPackageName, packageVersion)
	cmd := exec.Command("nuget", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("NuGet install package failed: %v\n", err)
		return
	}

	// 找出安裝後的套件目錄
	var packageInstallDir string
	filepath.Walk(tempDir, func(path string, info os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}
		if info.IsDir() && strings.HasPrefix(filepath.Base(path), nugetPackageName+".") {
			packageInstallDir = path
			return filepath.SkipDir
		}
		return nil
	})

	if packageInstallDir == "" {
		fmt.Printf("Could not find installed package directory for %s\n", nugetPackageName)
		return
	}

	// 若未指定版本，從目錄名推斷版本
	if packageVersion == "" || packageVersion == "latest" {
		baseName := filepath.Base(packageInstallDir)
		prefix := nugetPackageName + "."
		if strings.HasPrefix(baseName, prefix) {
			packageVersion = strings.TrimPrefix(baseName, prefix)
		} else {
			packageVersion = "1.0.0"
		}
	}

	// 找可用的 framework
	libPath := filepath.Join(packageInstallDir, "lib")
	frameworkDirs := []string{}
	filepath.Walk(libPath, func(path string, info os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}
		if info.IsDir() && path != libPath {
			rel, _ := filepath.Rel(libPath, path)
			frameworkDirs = append(frameworkDirs, rel)
			return filepath.SkipDir
		}
		return nil
	})

	if len(frameworkDirs) == 0 {
		fmt.Printf("No target frameworks found under 'lib' for package %s.\n", nugetPackageName)
		return
	}

	// 自動選擇最適合 Unity 的framework
	selectedFramework := chooseFrameworkAuto(frameworkDirs)

	// 建立輸出目錄
	err = os.MkdirAll(filepath.Join(pluginPath, "Runtime"), os.ModePerm)
	if err != nil {
		fmt.Printf("Creation of the plugin directory failed: %v\n", err)
		return
	}

	fmt.Printf("Using target framework: %s\n", selectedFramework)
	frameworkDirPath := filepath.Join(libPath, selectedFramework)

	// 複製該框架目錄下的 dll
	dllFiles, err := filepath.Glob(filepath.Join(frameworkDirPath, "*.dll"))
	if err != nil {
		fmt.Printf("Error locating DLLs: %v\n", err)
		return
	}

	totalCopied := 0
	var asmName, dllName string
	if len(dllFiles) > 0 {
		dllName = filepath.Base(dllFiles[0])
		asmName = strings.TrimSuffix(dllName, filepath.Ext(dllName))
	}

	for _, dll := range dllFiles {
		dllNameLocal := filepath.Base(dll)
		destPath := filepath.Join(pluginPath, "Runtime", dllNameLocal)
		fmt.Printf("Copying: %s\n", dllNameLocal)
		err = copyFile(dll, destPath)
		if err != nil {
			fmt.Printf("Error copying %s: %v\n", dllNameLocal, err)
			return
		}
		totalCopied++
	}

	// 建立 package.json
	err = createPackageJson(nugetPackageName, packageVersion, pluginPath)
	if err != nil {
		fmt.Printf("Error creating package.json: %v\n", err)
		return
	}

	// 建立 asmdef (如果有 DLL)
	if asmName != "" && dllName != "" {
		err = createAsmdef(asmName, dllName, filepath.Join(pluginPath, "Runtime"))
		if err != nil {
			fmt.Printf("Error creating asmdef: %v\n", err)
			return
		}
		fmt.Printf("Created asmdef for: %s\n", asmName)
	}

	fmt.Printf("\n========== Script finished, copied [%d] DLL(s) from '%s' to %s! ==========\n", totalCopied, selectedFramework, pluginPath)
	fmt.Println("Now creating .unitypackage without using Unity...")

	unityPackageName := nugetPackageName + ".unitypackage"
	err = createUnityPackageFromExport(pluginPath, nugetPackageName, unityPackageName)
	if err != nil {
		fmt.Printf("Error creating unitypackage: %v\n", err)
		return
	}

	fmt.Printf("Unitypackage '%s' created successfully!\n", unityPackageName)
}
