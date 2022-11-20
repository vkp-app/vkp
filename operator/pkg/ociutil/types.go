package ociutil

import "github.com/google/go-containerregistry/pkg/authn"

type DockerConfigJSON struct {
	Auths map[string]authn.AuthConfig `json:"auths,omitempty"`
}
