# Go NuGet Unity Exporter

Go NuGet Unity Exporter is a simple yet powerful tool for extracting DLL files from NuGet packages and exporting them for use in Unity projects. This tool is written in Go and provides a straightforward command-line interface.

## Features

- Download and extract DLL files from NuGet packages
- Maintain original directory structure
- Cross-platform support (Windows, macOS, Linux)

## Prerequisites

- Go 1.16 or higher
- NuGet command-line tool (ensure it's in your system PATH)

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

## Project Structure

```
go-nuget-unity-exporter/
├── nuget_exporter.go
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## Contributing

Contributions are welcome! If you have any suggestions for improvements or bug reports, please feel free to create an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

This project was inspired by [NuGet_to_Unity](https://github.com/eucylin/NuGet_to_Unity). Thanks to eucylin for the original work.