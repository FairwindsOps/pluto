FROM golang:1.14.3 AS build-env
WORKDIR /go/src/github.com/fairwindsops/pluto/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build-linux


FROM alpine:3.12.0 as alpine
RUN apk --no-cache --update add ca-certificates tzdata && update-ca-certificates


FROM scratch
COPY --from=alpine /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine /etc/passwd /etc/passwd

USER nobody
COPY --from=build-env /go/src/github.com/fairwindsops/pluto /

WORKDIR /opt/app

CMD ["/pluto"]
