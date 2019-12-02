FROM golang:1.13-alpine

ENV APP_DIR $GOPATH/src/github.com/stefanoj3/prometheus-playground

RUN apk add --no-cache bash make git tar curl

COPY . $APP_DIR
WORKDIR $APP_DIR

RUN go mod download
RUN go build -o /bin/server cmd/server/main.go

EXPOSE 8099
ENTRYPOINT ["/bin/server"]