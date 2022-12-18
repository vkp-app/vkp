package cluster

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	vclusterv1alpha1 "github.com/loft-sh/cluster-api-provider-vcluster/api/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"os"
	"sigs.k8s.io/cluster-api/api/v1beta1"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"text/template"
)

//go:embed config/values.yaml
var valuesTemplate string

var valuesTpl = template.Must(template.New("values.yaml").Parse(valuesTemplate))

func VCluster(ctx context.Context, cluster *v1alpha1.Cluster, version *v1alpha1.ClusterVersion, dexCustomCA bool, customCA, haConnectionString string) (*vclusterv1alpha1.VCluster, error) {
	log := logging.FromContext(ctx)
	hostname := getHostname(cluster)
	values := new(bytes.Buffer)
	// allow the platform to specify a storage class
	// if the user doesn't
	if cluster.Spec.Storage.StorageClassName == "" {
		cluster.Spec.Storage.StorageClassName = os.Getenv(EnvStorageClass)
	}
	envVars := map[string]string{}
	for _, kv := range os.Environ() {
		k, v, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		if strings.HasPrefix(k, "__VKP_") || strings.HasPrefix(kv, "__GLASS_") {
			log.V(5).Info("found environment variable to pass down to clusters", "key", k)
			envVars[k] = v
		}
	}
	// ensure that we have all the information we need to set up
	// an HA cluster
	enableHA := cluster.Spec.HA.Enabled && haConnectionString != ""
	valuesConfig := ValuesTemplate{
		Name: cluster.GetName(),
		Ingress: ValuesIngress{
			Host:          strings.TrimPrefix(hostname, "api."),
			TLSSecretName: IngressSecretName(cluster.GetName()),
			ClassName:     getEnv(EnvIngressClass, "nginx"),
			Issuer:        getEnv(EnvIngressIssuer, ""),
		},
		IDP: ValuesIDP{
			URL:        getEnv(EnvIDPURL, ""),
			SecretName: DexSecretName(cluster.GetName()),
			CustomCA:   dexCustomCA,
		},
		Storage: cluster.Spec.Storage,
		HA: ValuesHA{
			Enabled:      enableHA,
			Connection:   fmt.Sprintf("%s?sslmode=require", strings.ReplaceAll(haConnectionString, "postgresql://", "postgres://")),
			ReplicaCount: 2,
		},
		OpenShift:     getEnv(EnvIsOpenShift, "false") == "true",
		Image:         version.Spec.Image.String(),
		VclusterImage: os.Getenv(EnvVclusterImage),
		CoreDNSImage:  getEnv(EnvCoreDNSImage, "coredns/coredns:1.8.7"),
		CustomCA:      customCA,
		Plugins: ValuesPlugins{
			SyncImage:  getEnv(EnvSyncImage, "dev.local/vkp/vcluster-plugin-sync:latest"),
			HookImage:  getEnv(EnvHookImage, "dev.local/vkp/vcluster-plugin-hooks:latest"),
			PullPolicy: getEnv(EnvPluginPolicy, "Never"),
		},
		EnvVars: envVars,
	}
	log.V(3).Info("templating values.yaml file", "Template", valuesTemplate, "Overrides", valuesConfig)
	if err := valuesTpl.Execute(values, valuesConfig); err != nil {
		log.Error(err, "failed to template values.yaml")
		return nil, err
	}
	return &vclusterv1alpha1.VCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetName(),
			Namespace: cluster.GetNamespace(),
			Labels:    Labels(cluster),
		},
		Spec: vclusterv1alpha1.VClusterSpec{
			ControlPlaneEndpoint: v1beta1.APIEndpoint{
				Host: hostname,
				Port: 443,
			},
			HelmRelease: &vclusterv1alpha1.VirtualClusterHelmRelease{
				Chart: vclusterv1alpha1.VirtualClusterHelmChart{
					Name:    getOrDefault(version.Spec.Chart.Name, getEnv(EnvChartName, "vcluster")),
					Repo:    getOrDefault(version.Spec.Chart.Repository, getEnv(EnvChartRepo, "https://charts.loft.sh")),
					Version: getOrDefault(version.Spec.Chart.Version, getEnv(EnvChartVersion, "0.13.0")),
				},
				Values: values.String(),
			},
			// this will be ignored since we're manually setting the image version
			// above
			KubernetesVersion: pointer.String(getEnv(EnvKubeVersion, "1.25")),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// getOrDefault returns the first argument if
// it is a non-empty string. Otherwise, it returns
// the second argument.
func getOrDefault(v1, v2 string) string {
	if v1 != "" {
		return v1
	}
	return v2
}
