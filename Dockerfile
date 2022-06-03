FROM alpine:3.16.0

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
