# Stage 1: Build the Go binary
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go.mod
COPY go.mod ./

# 下載依賴
RUN go mod download

# 複製其餘的源代碼
COPY . .

# 構建靜態的 Go 服務器二進制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server ./cmd/server

# Stage 2: Runtime stage 基于 Debian slim
FROM debian:bullseye-slim AS runtime

WORKDIR /app

# 安装 Mono 和其他必要的包
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        mono-complete \
        wget \
        ca-certificates \
        bash && \
    rm -rf /var/lib/apt/lists/*

# 下载最新的 nuget.exe
RUN wget -O /app/nuget.exe https://dist.nuget.org/win-x86-commandline/latest/nuget.exe

# 创建一个包装 nuget.exe 的脚本并将其加入到 PATH 中
RUN printf '#!/bin/sh\nmono /app/nuget.exe "$@"\n' > /usr/local/bin/nuget && \
    chmod +x /usr/local/bin/nuget

# 复制从 builder 阶段构建的 Go 二进制文件
COPY --from=builder /app/server /app/server

# 确保二进制文件具有执行权限
RUN chmod +x /app/server

# 设置环境变量
ENV HOME=/root
ENV NUGET_PACKAGES=/root/.nuget/packages
ENV PATH="/app:/usr/local/bin:${PATH}"

# 设置默认命令来运行 Go 服务器
CMD ["./server"]