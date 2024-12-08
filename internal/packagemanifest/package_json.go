package packagemanifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreatePackageJson(packageName, version, outputPath string) error {
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
