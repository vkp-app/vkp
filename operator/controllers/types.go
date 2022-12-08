package controllers

const (
	labelTenant = "paas.dcas.dev/tenant"
	labelOwned  = "paas.dcas.dev/owned"

	finalizer = "paas.dcas.dev/finalizer"

	secretKeyDbConn = "pgbouncer-uri"
)

type TenantOptions struct {
	SkipDefaultAddons bool
	CustomCAFile      string

	NamespaceOwnership bool
	NamespaceLabels    bool
}

type ClusterOptions struct {
	DexGrpcAddr               string
	AllowHA                   bool
	UseHANonce                bool
	PostgresResourceName      string
	PostgresResourceNamespace string
}
