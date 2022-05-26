FROM golang:1.18

ADD . /app
WORKDIR /app

RUN go build

ENTRYPOINT ["./mood-detector"]
