@echo off
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o nuget-exporter.exe nuget_exporter.go