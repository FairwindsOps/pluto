FROM alpine:3.17

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
