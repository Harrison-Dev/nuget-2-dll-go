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
	err := internal.ExportNugetPackageToUnity(packageName, packageVersion, exportDir)
	if err != nil {
		log.Printf("Error exporting package: %v\n", err)
		http.Error(w, "Failed to export unitypackage", http.StatusInternalServerError)
		return
	}

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

	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Error sending file: %v\n", err)
	}

	os.Remove(unityPackagePath)

}
