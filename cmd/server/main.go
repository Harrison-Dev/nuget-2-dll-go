package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Harrison-Dev/nuget-2-dll-go/internal"
)

func main() {
	http.HandleFunc("/download", downloadHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	packageName := r.URL.Query().Get("package_name")
	if packageName == "" {
		http.Error(w, "package_name is required", http.StatusBadRequest)
		return
	}

	packageVersion := r.URL.Query().Get("package_version")
	if packageVersion == "" {
		packageVersion = "latest"
	}

	exportDir := "./export"

	// 呼叫我們的內部邏輯
	err := internal.ExportNugetPackageToUnity(packageName, packageVersion, exportDir)
	if err != nil {
		log.Printf("Error exporting package: %v\n", err)
		http.Error(w, "Failed to export unitypackage", http.StatusInternalServerError)
		return
	}

	// unitypackage 名稱就是 packageName + ".unitypackage"
	unityPackageName := packageName + ".unitypackage"
	unityPackagePath := filepath.Join(".", unityPackageName)

	fileInfo, err := os.Stat(unityPackagePath)
	if os.IsNotExist(err) {
		http.Error(w, "unitypackage not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error accessing unitypackage: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileInfo.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	file, err := os.Open(unityPackagePath)
	if err != nil {
		http.Error(w, "Unable to open unitypackage", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = ioCopy(w, file)
	if err != nil {
		log.Printf("Error sending file: %v\n", err)
		// 回傳前端可能已中斷，不一定要視為錯誤
	}

	// 可根據需要清理產生的檔案，避免長期累積
	os.Remove(unityPackagePath)
}

func ioCopy(dst http.ResponseWriter, src *os.File) (int64, error) {
	return io.Copy(dst, src)
}
