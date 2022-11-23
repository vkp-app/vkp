package controllers

const (
	labelTenant = "paas.dcas.dev/tenant"
	finalizer   = "paas.dcas.dev/finalizer"
)

type TenantOptions struct {
	SkipDefaultAddons bool
	CustomCAFile      string
}

type ClusterOptions struct {
	DexGrpcAddr string
}
