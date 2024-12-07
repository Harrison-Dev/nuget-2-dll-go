package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func createPackageJson(packageName, version, outputPath string) error {
	// 將 packageName 標準化為 Unity 可接受的格式 (小寫、用 '-' 代替 '.')
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
	// 為了避免與 DLL 同名，給 asmdef 加一個後綴
	// name 欄位中則可用合適的命名。例如：將 asmName 做小寫處理。
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

func selectFromList(prompt string, items []string) int {
	fmt.Println(prompt)
	for i, item := range items {
		fmt.Printf("[%d] %s\n", i, item)
	}
	fmt.Print("Select an index: ")
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil || choice < 0 || choice >= len(items) {
		fmt.Println("Invalid selection.")
		return -1
	}
	return choice
}

func main() {
	fmt.Println("Welcome to the Interactive NuGet to Unity Package Exporter!")
	nugetPackageName := getUserInput("Enter the NuGet package name (e.g. Newtonsoft.Json)", "")
	packageVersion := getUserInput("Enter the package version (or leave empty for latest)", "")

	// 建立暫存目錄
	tempDir, err := os.MkdirTemp("", "nuget_temp")
	if err != nil {
		fmt.Printf("Failed to create temporary directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir) // 結束時清理

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
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, wErr error) error {
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

	// 如果使用者沒指定版本，嘗試從目錄名推斷版本
	if packageVersion == "" || packageVersion == "latest" {
		baseName := filepath.Base(packageInstallDir)
		prefix := nugetPackageName + "."
		if strings.HasPrefix(baseName, prefix) {
			packageVersion = strings.TrimPrefix(baseName, prefix)
		} else {
			// fallback預設
			packageVersion = "1.0.0"
		}
	}

	// 找出可用的 target frameworks
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

	choice := selectFromList("Select a target framework to use:", frameworkDirs)
	if choice < 0 {
		return
	}
	selectedFramework := frameworkDirs[choice]

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
	var asmName string
	var dllName string
	if len(dllFiles) > 0 {
		// 假設取第一個 DLL 的檔名作為 asmdef 名稱的基底
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
	fmt.Println("You can now use this folder as a Unity package by placing it under your project's 'Packages' directory, or reference it via a local package path.")
	fmt.Println("If you have a test assembly, ensure that its asmdef references the newly created asmdef to access the DLL's namespaces.")
}
