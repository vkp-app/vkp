package syncers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddonSyncer_manifestsFromOCI(t *testing.T) {
	as := &AddonSyncer{}

	path, err := as.manifestsFromOCI(context.TODO(), "harbor.dcas.dev/docker.io/k8slt/imgpkg-test")
	assert.NoError(t, err)
	t.Log(path)
}
