

#FROM golang:alpine
#WORKDIR /app
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

FROM docker.io/jeanblanchard/alpine-glibc AS final

WORKDIR /app
COPY ./main                 /app/main
#COPY ./conf.yaml               /app/
COPY ./static              /app/static
COPY ./tpl             /app/tpl

# 清空缓存
RUN rm -rf /var/cache/apk/*

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8 \
    USER_HOME_DIR="/root"

ENTRYPOINT ["/app/main"]

