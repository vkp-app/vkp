#### Images

**Control plane**

`ghcr.io/vkp-app/vkp/apiserver:VERSION`

`ghcr.io/vkp-app/vkp/metrics-proxy:VERSION`

`ghcr.io/vkp-app/vkp/web:VERSION`

`ghcr.io/vpk-app/vkp/dex:VERSION`

**Data plane**

`ghcr.io/vkp-app/vkp/vcluster-plugin-hooks:VERSION`

`ghcr.io/vkp-app/vkp/vcluster-plugin-sync:VERSION`

**Operator/OLM**

`ghcr.io/vkp-app/vkp/operator:VERSION`

`ghcr.io/vkp-app/vkp/bundle:VERSION`<sup>1</sup>

`ghcr.io/vkp-app/vkp/index:VERSION`<sup>1</sup>

1. These images are only required if installing via the [Operator Lifecycle Manager](https://olm.operatorframework.io/) (e.g. on OpenShift).

**Helm charts**

> Helm charts are packaged as OCI artefacts.
> See the [documentation](https://helm.sh/docs/topics/registries/) for more details.

`ghcr.io/vkp-app/vkp/helm-charts`

For installation guides and more information; check the [documentation](https://vkp-app.github.io/docs/operator-guide/getting-started/).
