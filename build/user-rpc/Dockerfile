FROM golang:1.17-buster as builder

WORKDIR /app

# 换成阿里云镜像
RUN sed -i "s/archive.ubuntu./mirrors.aliyun./g" /etc/apt/sources.list
RUN sed -i "s/deb.debian.org/mirrors.aliyun.com/g" /etc/apt/sources.list
RUN sed -i "s/security.debian.org/mirrors.aliyun.com\/debian-security/g" /etc/apt/sources.list

# 安装依赖
COPY go.* ./
# 设置代理
RUN go env -w GOPROXY="https://goproxy.cn,direct"
# 下载依赖
RUN go mod tidy

# 将代码文件写入镜像
COPY . ./

# 构建二进制文件
RUN go build -mod=readonly -v -o user-server services/user/rpc/main.go

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# 将构建好的文件拷贝进镜像
COPY --from=builder /app/user-server /app/user-server

# 启动 Web 服务
CMD ["/app/user-server"]