FROM alpine:3.16.3

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
