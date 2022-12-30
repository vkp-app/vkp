package clusterutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewID(t *testing.T) {
	id := NewID()
	t.Log(id)
	assert.Len(t, id, idLength)
}
