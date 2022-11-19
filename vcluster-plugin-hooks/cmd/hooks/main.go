package main

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-hooks/internal/hooks"
	"os"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

const (
	EnvClusterDomain = "VCLUSTER_CLUSTER_DOMAIN"
	EnvClusterName   = "VCLUSTER_CLUSTER_NAME"
)

func init() {
	// initialise all the APIs we're going to use
	_ = gatewayv1beta1.AddToScheme(plugin.Scheme)
}

func main() {
	clusterDomain := os.Getenv(EnvClusterDomain)
	clusterName := os.Getenv(EnvClusterName)

	// start the plugin
	_ = plugin.MustInit()
	plugin.MustRegister(&hooks.IngressHook{ClusterDomain: clusterDomain})
	plugin.MustRegister(&hooks.GatewayHttpHook{ClusterDomain: clusterDomain})
	plugin.MustRegister(&hooks.GatewayHook{ClusterDomain: clusterDomain})
	plugin.MustRegister(&hooks.PodHook{ClusterName: clusterName})
	plugin.MustStart()
}
