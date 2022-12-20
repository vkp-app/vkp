package main

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/syncers"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/syncers/crdresources"
	"os"
)

func main() {
	clusterName := os.Getenv(syncers.EnvClusterName)
	namespace := os.Getenv(syncers.EnvNamespace)

	// init the plugin
	ctx := plugin.MustInit()
	plugin.MustRegister(syncers.NewRBACSyncer())
	plugin.MustRegister(syncers.NewAddonSyncer(clusterName))
	plugin.MustRegister(syncers.NewSecretSyncer(clusterName, namespace))
	plugin.MustRegister(syncers.NewCertificateSyncer(clusterName))

	// resource sync
	plugin.MustRegister(crdresources.NewServiceMonitorSyncer(ctx))

	// start
	plugin.MustStart()
}
