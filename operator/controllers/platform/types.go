package platform

const annotationMainContainer = "kubectl.kubernetes.io/default-container"
const annotationConfigHash = "paas.dcas.dev/config-hash"

const SecretKeyOauthCookie = "cookieSecret"
const SecretKeyDexClientSecret = "dexClientSecret"

const (
	ComponentApiServer  = "apiserver"
	ComponentOauthProxy = "oauth-proxy"

	ComponentWeb = "web"
	ComponentDex = "dex"
)
