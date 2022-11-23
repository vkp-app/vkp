# vkp

![Version: 0.2.3](https://img.shields.io/badge/Version-0.2.3-informational?style=flat-square) ![AppVersion: 0.1.0](https://img.shields.io/badge/AppVersion-0.1.0-informational?style=flat-square)

## Values

| Key                             | Type   | Default                                               | Description                                                                                                                                            |
|---------------------------------|--------|-------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| api.image                       | string | `"ghcr.io/vkp-app/vkp/apiserver:main"`                | APIServer image                                                                                                                                        |
| api.imagePullPolicy             | string | `"IfNotPresent"`                                      | APIServer image pull policy                                                                                                                            |
| dex.image                       | string | `"ghcr.io/dexidp/dex:v2.35.3"`                        | Dex image                                                                                                                                              |
| dex.imagePullPolicy             | string | `"IfNotPresent"`                                      | Dex image pull policy                                                                                                                                  |
| dex.ingress.tlsSecret           | string | `"tls-dex"`                                           | Dex TLS certificate                                                                                                                                    |
| global.caSecret                 | string | `""`                                                  | Custom Certificate Authority to use for all components. Generally this should contain a single CA, but it can support many.                            |
| global.ingress.annotations      | object | `{}`                                                  | Annotations to add to all ingress resources (e.g. cert-manager issuers)                                                                                |
| global.ingress.domain           | string | `"example.org"`                                       | Base domain for components to be hosted on.                                                                                                            |
| global.ingress.ingressClassName | string | `""`                                                  | IngressClass that will be used to implement the Ingress.                                                                                               |
| global.ingress.tlsSecret        | string | `"tls-vkp"`                                           | TLS certificate to use for ingress. Doesn't need to exist if cert-manager is creating it (assuming you have set your annotations correctly).           |
| idp.clientSecret                | string | `""`                                                  | OIDC client secret that VKP components will use to authenticate to Dex                                                                                 |
| idp.connectors                  | list   | `[{"id":"mock","name":"Mock","type":"mockCallback"}]` | Dex connectors that VKP will delegate authentication to. https://dexidp.io/docs/connectors/                                                            |
| idp.cookieSecret                | string | `""`                                                  | Secret to use for the Oauth proxy cookies                                                                                                              |
| idp.existingSecret              | string | `""`                                                  | Existing secret to load credentials from. Must contain DEX_CLIENT_SECRET and OAUTH2_PROXY_COOKIE_SECRET keys                                           |
| oauthProxy.image                | string | `"quay.io/oauth2-proxy/oauth2-proxy:v7.4.0-amd64"`    | Oauth proxy image                                                                                                                                      |
| oauthProxy.imagePullPolicy      | string | `"IfNotPresent"`                                      | Oauth proxy image pull policy                                                                                                                          |
| prometheus.extraMetrics         | list   | `[]`                                                  | Additional metrics to show on the cluster overview page.                                                                                               |
| prometheus.url                  | string | `""`                                                  | Url to the prometheus instance. Embedded environment variables will be expanded e.g. `http://$PROMETHEUS_USERNAME:$PROMETHEUS_PASSWORD@promtheus:9090` |
| web.image                       | string | `"ghcr.io/vkp-app/vkp/web:main"`                      | Web image                                                                                                                                              |
| web.imagePullPolicy             | string | `"IfNotPresent"`                                      | Web image pull policy                                                                                                                                  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
