package syncers

import (
	"context"
	"github.com/stretchr/testify/assert"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	"testing"
)

func TestAddonSyncer_manifestsFromOCI(t *testing.T) {
	as := &AddonSyncer{}

	path, err := as.manifestsFromOCI(context.TODO(), paasv1alpha1.OCIRemoteRef{
		Name: "harbor.dcas.dev/docker.io/k8slt/imgpkg-test",
	})
	assert.NoError(t, err)
	t.Log(path)
}
