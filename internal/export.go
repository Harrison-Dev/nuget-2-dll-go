package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Harrison-Dev/nuget-2-dll-go/internal/nuget"
	"github.com/Harrison-Dev/nuget-2-dll-go/internal/packagemanifest"
	"github.com/Harrison-Dev/nuget-2-dll-go/internal/unitypackage"
)

// ExportNugetPackageToUnity 是高階函式，整合所有功能：
// 1. 使用 nuget 下載指定套件
// 2. 選擇框架並複製 DLL
// 3. 建立 package.json 與 asmdef
// 4. 將結果打包成 unitypackage
func ExportNugetPackageToUnity(nugetPackageName, packageVersion, exportPath string) error {
	// 建立暫存目錄
	tempDir, err := os.MkdirTemp("", "nuget_temp")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	pluginPath := filepath.Join(exportPath, nugetPackageName)

	// 下載套件
	fmt.Printf("\nDownloading %s (%s) to temporary directory...\n", nugetPackageName, packageVersion)
	err = nuget.DownloadPackage(nugetPackageName, packageVersion, tempDir)
	if err != nil {
		return fmt.Errorf("NuGet install package failed: %v", err)
	}

	// 找到安裝目錄
	packageInstallDir, err := nuget.FindInstalledPackageDir(nugetPackageName, tempDir)
	if err != nil {
		return err
	}

	// 從目錄推斷版本
	packageVersion, err = nuget.DetermineVersionIfEmpty(nugetPackageName, packageVersion, packageInstallDir)
	if err != nil {
		return err
	}

	// 找框架
	frameworkDirs, err := nuget.ListFrameworks(packageInstallDir)
	if err != nil {
		return err
	}

	selectedFramework := nuget.ChooseFrameworkAuto(frameworkDirs)
	fmt.Printf("Using target framework: %s\n", selectedFramework)

	// 複製 DLL
	err = os.MkdirAll(filepath.Join(pluginPath, "Runtime"), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Creation of the plugin directory failed: %v", err)
	}

	dllName, asmName, totalCopied, err := nuget.CopyDlls(packageInstallDir, selectedFramework, filepath.Join(pluginPath, "Runtime"))
	if err != nil {
		return err
	}

	// 建立 package.json
	err = packagemanifest.CreatePackageJson(nugetPackageName, packageVersion, pluginPath)
	if err != nil {
		return fmt.Errorf("Error creating package.json: %v", err)
	}

	// 建立 asmdef
	if asmName != "" && dllName != "" {
		err = packagemanifest.CreateAsmdef(asmName, dllName, filepath.Join(pluginPath, "Runtime"))
		if err != nil {
			return fmt.Errorf("Error creating asmdef: %v", err)
		}
		fmt.Printf("Created asmdef for: %s\n", asmName)
	}

	fmt.Printf("\n========== Script finished, copied [%d] DLL(s) from '%s' to %s! ==========\n", totalCopied, selectedFramework, pluginPath)
	fmt.Println("Now creating .unitypackage without using Unity...")

	unityPackageName := nugetPackageName + ".unitypackage"
	err = unitypackage.CreateUnityPackageFromExport(pluginPath, nugetPackageName, unityPackageName)
	if err != nil {
		return fmt.Errorf("Error creating unitypackage: %v", err)
	}

	fmt.Printf("Unitypackage '%s' created successfully!\n", unityPackageName)
	return nil
}
