FROM ubuntu:latest
LABEL authors="levi.tang"

RUN rm /bin/sh && ln -s /bin/bash /bin/sh
# 安装程序
RUN apt update && apt install -y wget git
# 输入DNS服务器
RUN wget https://studygolang.com/dl/golang/go1.21.6.linux-amd64.tar.gz -O go1.21.6.linux-amd64.tar.gz && tar -zxvf ./go1.21.6.linux-amd64.tar.gz -C /usr/local
RUN echo "export PATH=$PATH:/usr/local/go/bin" >> /root/.bashrc && source /root/.bashrc && /usr/local/go/bin/go install github.com/go-delve/delve/cmd/dlv@latest

# 设置GO proxy
RUN /usr/local/go/bin/go env -w GOPROXY=https://goproxy.cn,direct && \
    /usr/local/go/bin/go env -w GO111MODULE=on && \
    /usr/local/go/bin/go  env -w GOPRIVATE=gitlab.vrviu.com/inviu
