FROM golang:1.14.3 AS build-env
WORKDIR /go/src/github.com/fairwindsops/pluto/

ARG version=dev
ARG commit=none

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go get github.com/markbates/pkger/cmd/pkger
RUN VERSION=$version COMMIT=$commit make build-linux


FROM alpine:3.12.0 as alpine
RUN apk --no-cache --update add ca-certificates tzdata && update-ca-certificates



USER nobody
COPY --from=build-env /go/src/github.com/fairwindsops/pluto /

WORKDIR /opt/app

ENTRYPOINT ["/pluto"]
