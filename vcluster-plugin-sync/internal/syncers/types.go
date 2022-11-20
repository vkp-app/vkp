package syncers

import "github.com/google/go-containerregistry/pkg/authn"

const (
	MagicClusterURL      = "__GLASS_CLUSTER_URL__"
	MagicDexURL          = "__GLASS_DEX_URL__"
	MagicDexClientID     = "__GLASS_DEX_CLIENT_ID__"
	MagicDexClientSecret = "__GLASS_DEX_CLIENT_SECRET__"
	MagicIngressClass    = "__GLASS_INGRESS_CLASS__"

	EnvClusterName = "VCLUSTER_CLUSTER_NAME"
)

const finalizer = "addon.paas.dcas.dev"

type DockerConfigJSON struct {
	Auths map[string]authn.AuthConfig `json:"auths,omitempty"`
}
