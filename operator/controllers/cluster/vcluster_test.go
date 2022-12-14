package cluster

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func Test_getOrDefault(t *testing.T) {
	var cases = []struct {
		v1  string
		v2  string
		out string
	}{
		{
			"foo",
			"bar",
			"foo",
		},
		{
			"",
			"bar",
			"bar",
		},
	}
	for _, tt := range cases {
		t.Run(tt.v1, func(t *testing.T) {
			out := getOrDefault(tt.v1, tt.v2)
			assert.EqualValues(t, tt.out, out)
		})
	}
}

func TestVCluster_chartConfig(t *testing.T) {
	cluster := &v1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.ClusterSpec{},
		Status: v1alpha1.ClusterStatus{
			ClusterID:     "1234",
			ClusterDomain: "example.org",
		},
	}

	t.Run("chart default values are used", func(t *testing.T) {
		vc, err := VCluster(context.TODO(), cluster, &v1alpha1.ClusterVersion{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ClusterVersionSpec{
				Image: v1alpha1.ClusterVersionImage{
					Repository: "rancher/k3s",
					Tag:        "v1.25.0-k3s.1",
				},
				Chart: v1alpha1.ClusterVersionChart{},
				Track: v1alpha1.TrackRegular,
			},
		}, false, "", "")
		assert.NoError(t, err)
		assert.EqualValues(t, "vcluster", vc.Spec.HelmRelease.Chart.Name)
	})
	t.Run("explicit chart overrides are used", func(t *testing.T) {
		vc, err := VCluster(context.TODO(), cluster, &v1alpha1.ClusterVersion{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.ClusterVersionSpec{
				Image: v1alpha1.ClusterVersionImage{
					Repository: "rancher/k3s",
					Tag:        "v1.25.0-k3s.1",
				},
				Chart: v1alpha1.ClusterVersionChart{
					Name: "vcluster-eks",
				},
				Track: v1alpha1.TrackRegular,
			},
		}, false, "", "")
		assert.NoError(t, err)
		assert.EqualValues(t, "vcluster-eks", vc.Spec.HelmRelease.Chart.Name)
	})
}
