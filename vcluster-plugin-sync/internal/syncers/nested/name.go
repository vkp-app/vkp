package nested

import (
	"encoding/base32"
	"hash/fnv"
	"strings"
)

// Kubernetes names must match the regexp '[a-z0-9]([-a-z0-9]*[a-z0-9])?'.
var encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

// idToName converts a resource name that may contain special characters into something
// that can be safely stored in Kubernetes.
//
// Borrowed from https://github.com/dexidp/dex/blob/master/storage/kubernetes/client.go
func idToName(s string) string {
	return strings.TrimRight(encoding.EncodeToString(fnv.New64().Sum([]byte(s))), "=")
}
