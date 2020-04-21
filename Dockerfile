FROM alpine:3.10

RUN apk add --no-cache curl

RUN curl -sSf https://installer.saltbo.cn/uptoc.sh | sh

COPY LICENSE README.md /
COPY scripts/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
