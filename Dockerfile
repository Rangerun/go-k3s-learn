# 第一阶段：编译阶段 builder
FROM golang:1.25-alpine AS builder

# ✅在这里增加国内Go代理配置！！！
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off
RUN apk add --no-cache ca-certificates

WORKDIR /build

# 先复制依赖文件，利用docker缓存层
COPY go.mod go.sum ./
RUN go mod download

# 复制全部源码
COPY . .

# CGO_ENABLED=0 静态编译，不依赖系统glibc；GOOS=linux
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# 第二阶段：运行阶段，只拿编译出来的二进制
FROM alpine:3.20

WORKDIR /app

# 从builder阶段拷贝编译好的二进制文件
COPY --from=builder /build/app ./

EXPOSE 8080
ENTRYPOINT ["./app"]