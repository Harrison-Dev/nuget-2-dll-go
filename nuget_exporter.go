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

func main() {
	fmt.Println("Welcome to the Interactive NuGet Package Exporter!")

	nugetPackageName := getUserInput("Enter the NuGet package name (e.g. Newtonsoft.Json)", "")

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "nuget_temp")
	if err != nil {
		fmt.Printf("Failed to create temporary directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir) // Clean up the temp directory at the end

	exportPath := "./export"

	fmt.Printf("\nDownloading %s to temporary directory...\n", nugetPackageName)
	cmd := exec.Command("nuget", "install", nugetPackageName, "-OutputDirectory", tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("NuGet install package failed: %v\n", err)
		return
	}

	err = os.MkdirAll(exportPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Creation of the export directory failed: %v\n", err)
		return
	}
	fmt.Printf("Successfully created the export directory\n")

	fmt.Println("\nCopying DLLs to export directory...")
	totalCopied := 0
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".dll") {
			relPath, err := filepath.Rel(tempDir, path)
			if err != nil {
				return err
			}
			destPath := filepath.Join(exportPath, relPath)
			destDir := filepath.Dir(destPath)
			err = os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				return err
			}
			fmt.Printf("Copying: %s\n", relPath)
			err = copyFile(path, destPath)
			if err != nil {
				return err
			}
			totalCopied++
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error while copying files: %v\n", err)
		return
	}

	fmt.Printf("\n==========Script finished, copied [%d] DLLs to %s!==========\n", totalCopied, exportPath)
}
