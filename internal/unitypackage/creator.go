package unitypackage

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"

	"github.com/Harrison-Dev/nuget-2-dll-go/internal/utils"
)

// CreateUnityPackageFromExport 掃描 export/<packageName> 下所有檔案，打包成 .unitypackage
func CreateUnityPackageFromExport(exportDir, packageName, outPackageName string) error {
	// 收集所有檔案
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
		unityPath := filepath.ToSlash(filepath.Join("Assets", packageName, rel))
		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		guid := utils.GenGUID()

		// 寫入 asset
		err = writeTarFile(tarWriter, guid+"/asset", content)
		if err != nil {
			return err
		}

		// 寫入 asset.meta
		metaContent := utils.GenerateMeta(guid)
		err = writeTarFile(tarWriter, guid+"/asset.meta", metaContent)
		if err != nil {
			return err
		}

		// 寫入 pathname
		err = writeTarFile(tarWriter, guid+"/pathname", []byte(unityPath))
		if err != nil {
			return err
		}
	}

	return nil
}

func writeTarFile(tw *tar.Writer, name string, data []byte) error {
	header := &tar.Header{
		Name: name,
		Mode: 0600,
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	_, err := tw.Write(data)
	return err
}
