FROM golang:1.12.1-alpine3.9 as builder

RUN set -x \
    && apk add --no-cache git

COPY . /rpcgen

WORKDIR /rpcgen

RUN set -x \
    && export CGO_ENABLED=0 \
    && go build -o /go/bin/rpc-generator ./main.go

# Executable image
FROM alpine:3.9
run set -x \
    && mkdir /var/test
COPY --from=builder /go/bin/rpc-generator /usr/local/sbin/rpc-generator
WORKDIR /var/test
