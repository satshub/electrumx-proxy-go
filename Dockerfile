##################################
#
# 构建可执行服务
#
##################################
FROM golang:1.21 as builder

ENV PROJECT electrumx-proxy

ADD ./  /data/go/src/$PROJECT

ENV GO111MODULE on

ENV GOPROXY https://goproxy.cn,direct

#如果目录不存在，则会自动创建
WORKDIR /data/go/src/$PROJECT

RUN apt-get update \
	&& apt-get install ca-certificates make gcc git -y \
    && go build



##################################
#
#准备完整的可执行环境，并启动服务
#
##################################
FROM ubuntu:18.04

ENV PROJECT electrumx-proxy

WORKDIR /usr/local/$PROJECT

RUN apt-get update  && apt-get install ca-certificates -y

COPY  --from=builder /data/go/src/$PROJECT/electrumx-proxy-go /usr/local/$PROJECT/

ADD ./config.toml /usr/local/$PROJECT/config.toml

CMD ./electrumx-proxy-go
