FROM golang:1.13

ENV APP_HOME /srv
WORKDIR $APP_HOME

COPY go.* $APP_HOME/
RUN go mod download

COPY . .
RUN make build && cp build/bin/uptoc /usr/local/bin

COPY scripts/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
