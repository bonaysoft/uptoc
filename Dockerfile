FROM golang:1.13

ENV APP_HOME /srv
WORKDIR $APP_HOME

ENV GOPROXY=https://goproxy.cn,direct
COPY go.* $APP_HOME/
RUN go mod download

COPY . .
RUN make build

ENTRYPOINT ["scripts/entrypoint.sh"]
