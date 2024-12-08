package nuget

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DownloadPackage 使用 nuget CLI 下載指定套件至 tempDir
func DownloadPackage(packageName, packageVersion, tempDir string) error {
	args := []string{"install", packageName, "-OutputDirectory", tempDir}
	if packageVersion != "" && packageVersion != "latest" {
		args = append(args, "-Version", packageVersion)
	}

	cmd := exec.Command("nuget", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// FindInstalledPackageDir 尋找安裝後的套件目錄
func FindInstalledPackageDir(packageName, tempDir string) (string, error) {
	var packageInstallDir string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}
		if info.IsDir() && strings.HasPrefix(filepath.Base(path), packageName+".") {
			packageInstallDir = path
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if packageInstallDir == "" {
		return "", fmt.Errorf("Could not find installed package directory for %s", packageName)
	}
	return packageInstallDir, nil
}
