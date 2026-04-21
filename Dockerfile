FROM alpine:3.23.4

LABEL org.opencontainers.image.authors="FairwindsOps, Inc." \
      org.opencontainers.image.vendor="FairwindsOps, Inc." \
      org.opencontainers.image.title="pluto" \
      org.opencontainers.image.description="Pluto is a cli tool to help discover deprecated apiVersions in Kubernetes" \
      org.opencontainers.image.documentation="https://pluto.docs.fairwinds.com/" \
      org.opencontainers.image.source="https://github.com/FairwindsOps/pluto" \
      org.opencontainers.image.url="https://github.com/FairwindsOps/pluto" \
      org.opencontainers.image.licenses="Apache License 2.0"

# Install CA bundle for TLS.
RUN apk --no-cache add ca-certificates

USER nobody
COPY pluto /

ENTRYPOINT ["/pluto"]
