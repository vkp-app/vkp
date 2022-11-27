package main

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/syncers"
	"os"
)

func main() {
	clusterName := os.Getenv(syncers.EnvClusterName)
	namespace := os.Getenv(syncers.EnvNamespace)

	// start the plugin
	_ = plugin.MustInit()
	plugin.MustRegister(syncers.NewRBACSyncer())
	plugin.MustRegister(syncers.NewAddonSyncer(clusterName))
	plugin.MustRegister(syncers.NewSecretSyncer(clusterName, namespace))
	plugin.MustStart()
}
