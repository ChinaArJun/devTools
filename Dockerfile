
# 项目docker环境
# base java shell openssl ca-certificates
# wkhtmltopdf 操作指令: wkhtmltopdf https://www.csdn.net /Users/kevin/Desktop/csdn.pdf
# imagemagick 操作指令：convert -density 200 -quality 100 test.pdf test.jpg

#FROM docker.io/jeanblanchard/alpine-glibc AS final
#FROM finebaas/base
FROM docker.io/jeanblanchard/alpine-glibc AS final
RUN apk --no-cache add ca-certificates wget \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub \
    https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub

MAINTAINER tools

WORKDIR /app
COPY .                 /app/


# 清空缓存
RUN rm -rf /var/cache/apk/*

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8 \
    USER_HOME_DIR="/root"

ENTRYPOINT ["/app/main"]

