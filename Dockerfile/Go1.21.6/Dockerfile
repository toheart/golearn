FROM ubuntu:latest
LABEL authors="levi.tang"

RUN rm /bin/sh && ln -s /bin/bash /bin/sh
# 安装程序
RUN apt update && apt install -y wget git build-essential zlib1g-dev libncurses5-dev libgdbm-dev libnss3-dev libssl-dev libreadline-dev libffi-dev wget
# 输入DNS服务器
RUN wget https://studygolang.com/dl/golang/go1.21.6.linux-amd64.tar.gz -O go1.21.6.linux-amd64.tar.gz && tar -zxvf ./go1.21.6.linux-amd64.tar.gz -C /usr/local
# 安装golang

# 设置GO proxy
RUN /usr/local/go/bin/go env -w GOPROXY=https://goproxy.cn,direct && \
    /usr/local/go/bin/go env -w GO111MODULE=on && \
    /usr/local/go/bin/go  env -w GOPRIVATE=gitlab.vrviu.com/inviu && \
     /usr/local/go/bin/go  env -w GOBIN=/root/go/bin  && /usr/local/go/bin/go install github.com/go-delve/delve/cmd/dlv@latest

# 安装python3
ADD ./Python-3.12.0.tgz /tmp
RUN cd /tmp/Python-3.12.0 && ./configure --enable-optimizations && make install

# 安装protoc相关组件
ADD ./protoc /tmp
RUN chmod a+x /tmp/protoc && mv /tmp/protoc /usr/local/bin && /usr/local/go/bin/go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && /usr/local/go/bin/go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest  \
    && /usr/local/go/bin/go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest \
    && /usr/local/go/bin/go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest \
    && /usr/local/go/bin/go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
    && /usr/local/go/bin/go install github.com/envoyproxy/protoc-gen-validate@latest \
    && /usr/local/go/bin/go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
    && /usr/local/go/bin/go install github.com/spf13/cobra-cli@latest

RUN apt install -y dnsutils iputils-ping curl iproute2  net-tools
