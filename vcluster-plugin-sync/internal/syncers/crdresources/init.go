package crdresources

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

const KoDataPathEnv = "KO_DATA_PATH"
const (
	CustomResourceServiceMonitorFile = "monitoring.coreos.com_servicemonitors.yaml"
)

func init() {
	_ = promv1.AddToScheme(plugin.Scheme)
	_ = promv1.AddToScheme(clientgoscheme.Scheme)
}
