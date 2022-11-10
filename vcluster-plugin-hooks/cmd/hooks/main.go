package main

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-hooks/internal/hooks"
	"os"
)

const EnvClusterDomain = "VCLUSTER_CLUSTER_DOMAIN"

func main() {
	clusterDomain := os.Getenv(EnvClusterDomain)

	// start the plugin
	_ = plugin.MustInit()
	plugin.MustRegister(&hooks.IngressHook{ClusterDomain: clusterDomain})
	plugin.MustStart()
}
