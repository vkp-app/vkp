package cluster

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	vclusterv1alpha1 "github.com/loft-sh/cluster-api-provider-vcluster/api/v1alpha1"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"os"
	"sigs.k8s.io/cluster-api/api/v1beta1"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
	"text/template"
)

//go:embed config/values.yaml
var valuesTemplate string

var valuesTpl = template.Must(template.New("values.yaml").Parse(valuesTemplate))

func VCluster(ctx context.Context, cluster *paasv1alpha1.Cluster) (*vclusterv1alpha1.VCluster, error) {
	log := logging.FromContext(ctx)
	hostname := fmt.Sprintf("%s.%s", cluster.Status.ClusterID, cluster.Status.ClusterDomain)
	values := new(bytes.Buffer)
	valuesConfig := ValuesTemplate{
		IngressClassName: getEnv(EnvIngressClass, "nginx"),
		Host:             hostname,
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
					Name:    getEnv(EnvChartName, "vcluster"),
					Repo:    getEnv(EnvChartRepo, "https://charts.loft.sh"),
					Version: getEnv(EnvChartVersion, "0.12.2"),
				},
				Values: values.String(),
			},
			// todo support release configuration (e.g. stable vs fast-track)
			KubernetesVersion: pointer.String(getEnv(EnvKubeVersion, "1.24")),
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
