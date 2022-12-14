# vkp

![Version: 0.6.1](https://img.shields.io/badge/Version-0.6.1-informational?style=flat-square) ![AppVersion: 0.0.0](https://img.shields.io/badge/AppVersion-0.0.0-informational?style=flat-square)

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| api.image.registry | string | `"ghcr.io"` | APIServer image registry |
| api.image.repository | string | `"vkp-app/vkp/apiserver"` | API Server image repository |
| api.image.tag | string | `""` | APIServer image tag (defaults to .Chart.AppVersion) |
| api.imagePullPolicy | string | `"IfNotPresent"` | APIServer image pull policy |
| dex.image.registry | string | `"ghcr.io"` | Dex image registry |
| dex.image.repository | string | `"vkp-app/vkp/dex"` | Dex image repository |
| dex.image.tag | string | `""` | Dex image tag |
| dex.imagePullPolicy | string | `"IfNotPresent"` | Dex image pull policy |
| dex.ingress.host | string | `"dex"` | Dex host. Always includes the `global.ingress.domain` as a suffix (e.g. setting this to "foo" will result in "foo.example.org"). |
| dex.ingress.tlsSecret | string | `"tls-dex"` | Dex TLS certificate |
| global.caSecret | string | `""` | Custom Certificate Authority to use for all components. Generally this should contain a single CA, but it can support many. |
| global.imagePullSecrets | list | `[]` | Global container registry secret names as an array. |
| global.imageRegistry | string | `""` | Global container image registry. Takes priority of any `image.registry` definitions. |
| global.ingress.annotations | object | `{}` | Annotations to add to all ingress resources (e.g. cert-manager issuers) |
| global.ingress.domain | string | `"example.org"` | Base domain for components to be hosted on. |
| global.ingress.ingressClassName | string | `""` | IngressClass that will be used to implement the Ingress. |
| global.ingress.tlsSecret | string | `"tls-vkp"` | TLS certificate to use for ingress. Doesn't need to exist if cert-manager is creating it (assuming you have set your annotations correctly). |
| idp.clientSecret | string | `""` | OIDC client secret that VKP components will use to authenticate to Dex |
| idp.connectors | list | `[{"id":"mock","name":"Mock","type":"mockCallback"}]` | Dex connectors that VKP will delegate authentication to. https://dexidp.io/docs/connectors/ |
| idp.cookieSecret | string | `""` | Secret to use for the Oauth proxy cookies |
| idp.existingSecret | string | `""` | Existing secret to load credentials from. Must contain `DEX_CLIENT_SECRET` and `OAUTH2_PROXY_COOKIE_SECRET` keys |
| metrics_proxy.enabled | bool | `true` | Enable or disable the Metrics Proxy component. |
| metrics_proxy.image.registry | string | `"ghcr.io"` | Metrics Proxy image registry |
| metrics_proxy.image.repository | string | `"vkp-app/vkp/metrics-proxy"` | Metrics Proxy image repository |
| metrics_proxy.image.tag | string | `""` | Metrics Proxy image tag |
| metrics_proxy.imagePullPolicy | string | `"IfNotPresent"` | Metrics Proxy image pull policy |
| oauthProxy.embedStaticResources | bool | `false` | Whether to use embedded static files (e.g. CSS). Required to work without an internet connection. |
| oauthProxy.image.registry | string | `"quay.io"` | Oauth proxy image registry |
| oauthProxy.image.repository | string | `"oauth2-proxy/oauth2-proxy"` | Oauth proxy image repository |
| oauthProxy.image.tag | string | `"v7.4.0-amd64"` | Oauth proxy image tag |
| oauthProxy.imagePullPolicy | string | `"IfNotPresent"` | Oauth proxy image pull policy |
| prometheus.extraMetrics | list | `[]` | Additional metrics to show on the cluster overview page. |
| prometheus.url | string | `""` | Url to the prometheus instance. Embedded environment variables will be expanded e.g. `http://$PROMETHEUS_USERNAME:$PROMETHEUS_PASSWORD@prometheus:9090` |
| vkp.certificates.certManagerNamespace | string | `"cert-manager"` | Name of the Cert Manager namespace that we can use to create ClusterIssuer's |
| vkp.certificates.issuer.create | bool | `true` | Create the self-signed issuer. If disabled, you will need to ensure that the issuer already exists. |
| vkp.certificates.issuer.kind | string | `"Issuer"` | Kind (Issuer/ClusterIssuer) of the resource to bootstrap from. |
| vkp.certificates.issuer.name | string | `"vkp-selfsigned"` | Name of the Issuer/ClusterIssuer to bootstrap from. |
| vkp.clusterVersions.default.enabled | bool | `false` | Install default ClusterVersions. Disable to supply your own. |
| vkp.clusterVersions.image.registry | string | `"docker.io"` | Container registry to pull images from. |
| vkp.clusterVersions.image.repository | string | `"rancher/k3s"` | Repository containing cluster images. |
| vkp.consoleHost | string | `"vkp"` | Console host. Always includes the `global.ingress.domain` as a suffix (e.g. setting this to "foo" will result in "foo.example.org"). |
| web.image.registry | string | `"ghcr.io"` | Web image registry |
| web.image.repository | string | `"vkp-app/vkp/web"` | Web image repository |
| web.image.tag | string | `""` | Web image tag (defaults to .Chart.AppVersion) |
| web.imagePullPolicy | string | `"IfNotPresent"` | Web image pull policy |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
