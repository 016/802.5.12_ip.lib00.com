# 使用官方 Golang 映像作为基础
FROM golang:1.23.4 AS builder

# 设置工作目录
WORKDIR /app

# 将当前目录的文件复制到工作目录
COPY . .

# 编译 Go 应用（指定输出文件名为 get_ip）
# 兼容mac os 添加命令 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 , 非mac os 可尝试删除
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dpit_get_ip .

# 第二阶段：使用更小的基础镜像
FROM alpine:3.21.0

# 设置工作目录
WORKDIR /root/

# 从 builder 阶段复制编译好的可执行文件
COPY --from=builder /app/dpit_get_ip .

RUN chmod +x dpit_get_ip

# 暴露应用的端口
EXPOSE 8080

# 指定容器启动时执行的命令
CMD ["./dpit_get_ip"]
