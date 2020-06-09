FROM golang:1.14-alpine as golang
WORKDIR /go/src/app/vendor/github.com/stojg/purr
COPY . .
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /go/bin/app

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
RUN zip -r -0 /zoneinfo.zip .

FROM scratch
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /go/bin/app /app
ENTRYPOINT ["/app"]
