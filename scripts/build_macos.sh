#!/bin/bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o nuget-exporter nuget_exporter.go