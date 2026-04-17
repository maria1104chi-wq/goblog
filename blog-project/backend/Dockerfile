# Go后端Dockerfile
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app/backend

# 安装必要工具
RUN apk add --no-cache git

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# 生产镜像
FROM alpine:latest

# 安装CA证书（用于HTTPS请求）
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建应用目录
WORKDIR /app/backend

# 从builder复制编译好的二进制文件
COPY --from=builder /app/backend/main .
COPY --from=builder /app/backend/static ./static

# 创建上传目录和日志目录
RUN mkdir -p /app/backend/static/uploads /app/backend/logs

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./main"]
