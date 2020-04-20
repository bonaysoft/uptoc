FROM alpine:3.10

RUN apk add --no-cache curl

RUN curl -sSf http://uptoc.saltbo.cn/install.sh | sh && cp $HOME/bin/uptoc /usr/local/bin/uptoc

COPY LICENSE README.md /

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
