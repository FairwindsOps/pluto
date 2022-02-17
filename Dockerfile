FROM alpine:3.15.0

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
