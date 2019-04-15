FROM golang:1.11.4-alpine
WORKDIR /go/src/craftli.co/reload/
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -a -o reload

FROM alpine:3.8
COPY --from=0 /go/src/craftli.co/reload/reload /usr/bin/reload
ENTRYPOINT ["reload"]
