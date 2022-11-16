FROM alpine:3.16.2

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
