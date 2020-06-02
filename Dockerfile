FROM golang:1.13 AS build-env
WORKDIR /go/src/github.com/fairwindsops/pluto/

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -a -o pluto *.go

FROM alpine:3.10
WORKDIR /usr/local/bin
RUN apk --no-cache add ca-certificates

RUN addgroup -S pluto && adduser -u 1200 -S pluto -G pluto
USER 1200
COPY --from=build-env /go/src/github.com/fairwindsops/pluto/pluto .

WORKDIR /opt/app

CMD ["pluto"]
