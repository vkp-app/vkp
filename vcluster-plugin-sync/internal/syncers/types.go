package syncers

import "github.com/google/go-containerregistry/pkg/authn"

const (
	EnvClusterName = "VCLUSTER_CLUSTER_NAME"
	EnvNamespace   = "KUBERNETES_NAMESPACE"
)

const finalizer = "addon.paas.dcas.dev"

type DockerConfigJSON struct {
	Auths map[string]authn.AuthConfig `json:"auths,omitempty"`
}
