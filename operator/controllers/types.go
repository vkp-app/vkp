package controllers

const (
	labelTenant = "paas.dcas.dev/tenant"
	finalizer   = "paas.dcas.dev/finalizer"
)

type TenantOptions struct {
	SkipDefaultAddons bool
}

type ClusterOptions struct {
	DexGrpcAddr string
}
