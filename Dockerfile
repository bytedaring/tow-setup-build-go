# --- 阶段 1: 构建 ---
# 使用一个完整的 Go + Alpine 镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /src

# 复制 go mod 文件并下载依赖，以便利用 Docker 的层缓存
# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

# 复制所有源代码
COPY . .

# 构建动态链接的 Go 应用
# Alpine 基础镜像默认 CGO 是启用的
RUN go build -o /app/main .

# --- 阶段 2: 运行 ---
# 使用极简的 alpine 作为运行环境
FROM alpine:latest

# 为操作系统安装最新的根证书，用于 TLS/SSL 验证
# 同时安装时区数据，很多应用都需要
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 只从构建阶段复制最终编译好的二进制文件
COPY --from=builder /app/main .

# （可选）如果你的服务需要访问特定文件，也从构建阶段复制
# COPY --from=builder /src/configs ./configs

# 暴露端口
EXPOSE 8080

# 最终运行的命令
CMD ["./main"]
