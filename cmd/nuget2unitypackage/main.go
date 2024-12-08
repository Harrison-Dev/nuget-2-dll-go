package main

import (
	"fmt"
	"os"

	"github.com/Harrison-Dev/nuget-2-dll-go/internal"
	"github.com/Harrison-Dev/nuget-2-dll-go/internal/utils"
)

func main() {
	fmt.Println("Welcome to the Interactive NuGet to Unity Package Exporter!")
	nugetPackageName := utils.GetUserInput("Enter the NuGet package name (e.g. Newtonsoft.Json)", "")
	packageVersion := utils.GetUserInput("Enter the package version (or leave empty for latest)", "")

	err := internal.ExportNugetPackageToUnity(nugetPackageName, packageVersion, "./export")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done.")
}
