# 第一階段：建構階段，基於 .NET SDK 映像
FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /app

# 安裝 Go
ENV GO_VERSION=1.20.4
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH=/usr/local/go/bin:$PATH

# 確認 Go 安裝成功
RUN go version

# 複製 Go 模組文件並下載依賴
COPY go.mod  ./
RUN go mod tidy
RUN go mod download

# 複製整個應用程式代碼
COPY . .

# 構建 Go 應用
RUN go build -o server ./cmd/server

# 第二階段：運行階段，基於 .NET 運行時映像
FROM mcr.microsoft.com/dotnet/runtime:7.0 AS runtime
WORKDIR /app

# 複製從建構階段編譯好的 Go 二進位檔
COPY --from=build /app/server /app/server

# 暴露應用所需的端口
EXPOSE 8080

# 設定容器啟動時運行的命令
CMD ["./server"]