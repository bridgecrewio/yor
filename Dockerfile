FROM golang:1.13-alpine

RUN apk update && apk add build-base git
WORKDIR ./yor
COPY "." "."
RUN go build

ENTRYPOINT ["./yor"]