package syncers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	"os"
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

func TestAddonSyncer_mangleYAML(t *testing.T) {
	require.NoError(t, os.Setenv("__VKP_FOO", "BAR"))
	require.NoError(t, os.Setenv("__GLASS_BAR", "FOO"))
	require.NoError(t, os.Setenv("_VKP_FOO", "ZOO"))

	test := `
foo: _VKP_FOO
zoo: __VKP_FOO
bar: __GLASS_BAR`

	r := new(AddonSyncer)

	out := r.mangleYAML(test)
	assert.EqualValues(t, out, `
foo: _VKP_FOO
zoo: BAR
bar: FOO`)
}
