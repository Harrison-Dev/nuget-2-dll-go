package nuget

import (
	"io"
	"os"
	"path/filepath"
	"sort"
)

// 預設的目標框架優先順序
var frameworkPriority = []string{
	"netstandard2.0",
	"net45",
	"net46",
	"net47",
	"net48",
	"netstandard2.1",
}

// ListFrameworks 列出該套件lib下的可用Framework
func ListFrameworks(packageInstallDir string) ([]string, error) {
	libPath := filepath.Join(packageInstallDir, "lib")
	frameworkDirs := []string{}
	err := filepath.Walk(libPath, func(path string, info os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}
		if info.IsDir() && path != libPath {
			rel, _ := filepath.Rel(libPath, path)
			frameworkDirs = append(frameworkDirs, rel)
			return filepath.SkipDir
		}
		return nil
	})
	return frameworkDirs, err
}

// ChooseFrameworkAuto 從找到的frameworkDirs中自動選擇最合適的
func ChooseFrameworkAuto(frameworkDirs []string) string {
	fwSet := make(map[string]bool)
	for _, fw := range frameworkDirs {
		fwSet[fw] = true
	}
	for _, preferred := range frameworkPriority {
		if fwSet[preferred] {
			return preferred
		}
	}
	sort.Strings(frameworkDirs)
	return frameworkDirs[0]
}

// CopyDlls 複製目標框架下的DLLs至指定路徑
func CopyDlls(packageInstallDir, selectedFramework, destPath string) (dllName, asmName string, totalCopied int, err error) {
	frameworkDirPath := filepath.Join(packageInstallDir, "lib", selectedFramework)
	dllFiles, err := filepath.Glob(filepath.Join(frameworkDirPath, "*.dll"))
	if err != nil {
		return "", "", 0, err
	}

	if len(dllFiles) > 0 {
		dllName = filepath.Base(dllFiles[0])
		asmName = dllName[:len(dllName)-len(filepath.Ext(dllName))]
	}

	for _, dll := range dllFiles {
		dllNameLocal := filepath.Base(dll)
		dst := filepath.Join(destPath, dllNameLocal)
		if copyErr := copyFile(dll, dst); copyErr != nil {
			return "", "", totalCopied, copyErr
		}
		totalCopied++
	}

	return dllName, asmName, totalCopied, nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 使用 io.Copy 取代自訂的 ioCopy
	_, err = io.Copy(destFile, sourceFile)
	return err
}

// 簡化 io.Copy
func ioCopy(dst, src *os.File) (int64, error) {
	buf := make([]byte, 32*1024)
	var written int64
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, ioErrShortWrite
			}
		}
		if er == nil {
			continue
		}
		if er == ioEOF {
			break
		}
		if er != nil {
			return written, er
		}
	}
	return written, nil
}

var ioEOF = EOFType{}
var ioErrShortWrite = ShortWriteErr{}

type EOFType struct{}

func (EOFType) Error() string { return "EOF" }

type ShortWriteErr struct{}

func (ShortWriteErr) Error() string { return "short write" }
