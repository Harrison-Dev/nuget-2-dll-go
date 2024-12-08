package packagemanifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CreateAsmdef(asmName, dllName, outputPath string) error {
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
