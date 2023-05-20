package controllers

const (
	finalizer       = "paas.dcas.dev/finalizer"
	secretKeyDbConn = "pgbouncer-uri"
)

const (
	ConditionAddonDigest         = "DigestResolved"
	ConditionAddonDigestResolved = "Resolved"
	ConditionAddonDigestErr      = "Error"

	ConditionVersion         = "VersionResolved"
	ConditionVersionResolved = "Resolved"
	ConditionVersionErr      = "Error"
)

type TenantOptions struct {
	SkipDefaultAddons bool
	CustomCAFile      string

	NamespaceOwnership bool
	NamespaceLabels    bool
}

type ClusterOptions struct {
	AllowHA                   bool
	UseHANonce                bool
	PostgresResourceName      string
	PostgresResourceNamespace string
	RootCAIssuerName          string
	RootCAIssuerKind          string
}
