package nuget

import (
	"path/filepath"
	"strings"
)

// DetermineVersionIfEmpty 若原本版本为空或為latest，從目錄名稱推斷版本
func DetermineVersionIfEmpty(packageName, packageVersion, packageInstallDir string) (string, error) {
	if packageVersion == "" || packageVersion == "latest" {
		baseName := filepath.Base(packageInstallDir)
		prefix := packageName + "."
		if strings.HasPrefix(baseName, prefix) {
			return strings.TrimPrefix(baseName, prefix), nil
		} else {
			return "1.0.0", nil
		}
	}
	return packageVersion, nil
}
