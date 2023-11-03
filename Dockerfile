FROM golang:1.21-alpine AS builder

RUN  go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /build
ADD go.mod .
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -trimpath -ldflags "-w -s" -o tcpdumpc main.go


FROM alpine:3.18
RUN apk update && apk add tcpdump

WORKDIR /root
COPY --from=builder /build/tcpdumpc /usr/local/bin/tcpdumpc

CMD ["tcpdumpc"]
