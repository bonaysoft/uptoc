FROM alpine:3.10

RUN apk add --no-cache curl

RUN curl -sSf https://uptoc.saltbo.cn/install.sh | sh

COPY LICENSE README.md /
COPY scripts/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
