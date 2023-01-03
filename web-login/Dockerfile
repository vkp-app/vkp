ARG DEX_VERSION=v2.35.3
FROM harbor.dcas.dev/ghcr.io/dexidp/dex:${DEX_VERSION}

ENV DEX_FRONTEND_DIR=/srv/dex/web

COPY --chown=1001:0 . /srv/dex/web