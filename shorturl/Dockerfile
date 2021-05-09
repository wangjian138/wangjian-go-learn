FROM golang:1.16.3
MAINTAINER "WANGJIAN"
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

COPY ./ /usr/local/shorturl/

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
go env -w GO111MODULE=on

EXPOSE 8888
EXPOSE 8080
ENTRYPOINT ["entrypoint.sh"]