package clusterutil

import (
	"math/rand"
	"strings"
)

const idRange = "abcdefghijklmnopqrstuvwxyz"
const idLength = 8

func NewID() string {
	var sb strings.Builder
	max := len(idRange)
	for i := 0; i < idLength; i++ {
		sb.WriteByte(idRange[rand.Intn(max)])
	}
	return sb.String()
}
