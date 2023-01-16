FROM alpine:3.17.1

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
