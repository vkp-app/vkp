ARG DEX_VERSION=v2.35.3
FROM harbor.dcas.dev/ghcr.io/dexidp/dex:${DEX_VERSION}

ENV DEX_FRONTEND_DIR=/srv/dex/web
LABEL org.opencontainers.image.source="https://github.com/vkp-app/vkp"

COPY --chown=1001:0 . /srv/dex/web