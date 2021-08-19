FROM alpine:3.14

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
