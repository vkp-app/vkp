package controllers

const (
	labelTenant = "paas.dcas.dev/tenant"
	finalizer   = "paas.dcas.dev/finalizer"

	secretKeyDbConn = "pgbouncer-uri"
)

type TenantOptions struct {
	SkipDefaultAddons  bool
	CustomCAFile       string
	NamespaceOwnership bool
}

type ClusterOptions struct {
	DexGrpcAddr               string
	AllowHA                   bool
	PostgresResourceName      string
	PostgresResourceNamespace string
}
