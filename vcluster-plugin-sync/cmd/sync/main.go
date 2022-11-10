package main

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	"gitlab.dcas.dev/k8s/kube-glass/vcluster-plugin-sync/internal/syncers"
)

func main() {
	_ = plugin.MustInit()
	plugin.MustRegister(syncers.NewRBACSyncer())
	plugin.MustRegister(syncers.NewAddonSyncer())
	plugin.MustStart()
}
