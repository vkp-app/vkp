# Developing

## Installation

Install `clusterctl`:

```bash
curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.2.5/clusterctl-linux-amd64 -o /tmp/clusterctl
install /tmp/clusterctl ~/.local/bin/clusterctl
```

Install resources:

```shell
# when using a fresh minikube
./firstrun.sh
# otherwise
make run
```

Install the APIServer:

```shell
cd apiserver/
kubectl apply -k k8s/
skaffold run
```

Install the Operator:

```shell
cd operator/
skaffold run
```

Install the Web proxy:

```shell
cd web/
skaffold run
```

Install the Metrics proxy:

```shell
cd metrics-proxy/
skaffold run
```
