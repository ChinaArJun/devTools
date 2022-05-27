
# 项目docker环境
# base java shell openssl ca-certificates
# wkhtmltopdf 操作指令: wkhtmltopdf https://www.csdn.net /Users/kevin/Desktop/csdn.pdf
# imagemagick 操作指令：convert -density 200 -quality 100 test.pdf test.jpg

FROM golang:alpine
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

FROM docker.io/jeanblanchard/alpine-glibc AS final

WORKDIR /app
COPY ./main                 /app/main
COPY ./conf.yaml               /app/
COPY ./static              /app/static
COPY ./tpl             /app/tpl

# 清空缓存
RUN rm -rf /var/cache/apk/*

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8 \
    USER_HOME_DIR="/root"

ENTRYPOINT ["/app/main"]

