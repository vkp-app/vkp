FROM quay.io/jetstack/cert-manager-controller:v1.8.0

WORKDIR /ko-app
COPY --chown=1000:0 main /ko-app/operator

ENTRYPOINT ["/ko-app/operator"]
