# Go NuGet Unity Exporter

Go NuGet Unity Exporter is a simple yet powerful tool for extracting DLL files from NuGet packages and exporting them for use in Unity projects. This tool is written in Go and provides a straightforward command-line interface.

## Features

- Download and extract DLL files from NuGet packages
- Maintain original directory structure
- Cross-platform support (Windows, macOS, Linux)
- Automated builds for Windows and macOS using GitHub Actions

## Prerequisites

Before using the Go NuGet Unity Exporter, ensure you have the following installed on your system:

1. **Go**: This tool is written in Go. Install Go from the [official Go website](https://golang.org/doc/install). Version 1.16 or higher is recommended.

2. **NuGet CLI**: The NuGet command-line interface is required for downloading packages. 

   - For Windows: 
     - Download the latest NuGet.exe from the [official NuGet website](https://www.nuget.org/downloads).
     - Add the directory containing NuGet.exe to your system's PATH.

   - For macOS/Linux:
     - Install Mono from the [official Mono website](https://www.mono-project.com/download/stable/).
     - Install NuGet using Mono: `sudo mono nuget.exe update -self`

Ensure all these tools are properly installed and accessible from your command line before proceeding with the installation and usage of Go NuGet Unity Exporter.

### Verifying Prerequisites

You can verify the installation of these tools by running the following commands in your terminal:

```bash
go version
nuget help
git --version
```

If any of these commands fail, please revisit the installation steps for the respective tool.

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/your-username/go-nuget-unity-exporter.git
   ```

2. Navigate to the project directory:
   ```
   cd go-nuget-unity-exporter
   ```

3. Build the executable:
   ```
   go build -o nuget-exporter
   ```

## Usage

1. Run the executable:
   ```
   ./nuget-exporter
   ```
   On Windows, use `nuget-exporter.exe`

2. Follow the prompts to enter the NuGet package name.

3. The program will download the NuGet package, extract the DLL files, and export them to the `./export` directory.

## Releases

This project uses GitHub Actions to build and release new versions. The release process is manual, allowing for version control and flexibility.

To download the latest version, visit the [Releases](https://github.com/your-username/go-nuget-unity-exporter/releases) page and download the appropriate file for your operating system.

## Troubleshooting

If you encounter any issues while using the Go NuGet Unity Exporter, try the following:

1. **NuGet command not found**: Ensure that NuGet is properly installed and added to your system's PATH.

2. **Permission denied errors**: Make sure you have the necessary permissions to write to the output directory.

3. **Unable to download packages**: Check your internet connection and firewall settings. Ensure you're not behind a proxy that's blocking NuGet.

4. **DLL not found in Unity**: Verify that the exported DLLs are placed in the correct directory within your Unity project.

If problems persist, please open an issue on the GitHub repository with detailed information about your environment and the error you're encountering.

## Contributing

Contributions are welcome! If you have any suggestions for improvements or bug reports, please feel free to create an issue or submit a pull request.

## Feedback and Contributions

Your feedback and contributions are welcome! If you have suggestions for improvements or encounter any issues, please feel free to:

- Open an issue in the GitHub repository
- Submit a pull request with your proposed changes
- Contact the maintainer directly

We appreciate your input in making this tool better for everyone.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

This project was inspired by [NuGet_to_Unity](https://github.com/eucylin/NuGet_to_Unity). Thanks to eucylin for the original work.