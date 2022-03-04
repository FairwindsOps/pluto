FROM alpine:3.15

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
