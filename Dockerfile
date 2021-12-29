FROM alpine:3.14

RUN apk add --no-cache git bash openssh
COPY pluto /

ENTRYPOINT ["/pluto"]
