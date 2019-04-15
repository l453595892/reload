FROM golang:1.11.4-alpine
WORKDIR /go/src/git.kunlun/KUNLUN-Hyper/k.rpc.resource
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -a -o k.rpc.resource

FROM alpine:3.8
ENV SRV_VERSION {{SRV_VERSION}}
COPY --from=0 /go/src/git.kunlun/KUNLUN-Hyper/k.rpc.resource/k.rpc.resource /usr/bin/k.rpc.resource
ENTRYPOINT ["k.rpc.resource"]