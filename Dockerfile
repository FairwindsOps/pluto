FROM alpine:3.14
RUN apk --no-cache --update add ca-certificates tzdata && update-ca-certificates

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
