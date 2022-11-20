package ociutil

import (
	"context"
	"encoding/json"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

func GetAuth(ctx context.Context, c client.Client, namespace string, oci *paasv1alpha1.OCIRemoteRef) ([]remote.Option, error) {
	log := logging.FromContext(ctx).WithValues("oci", oci.Name)
	auth := []remote.Option{remote.WithAuthFromKeychain(authn.DefaultKeychain)}
	if oci.ImagePullSecret.Name != "" {
		log.V(1).Info("fetching remote authentication information from secret", "secret", oci.ImagePullSecret)
		sec := &corev1.Secret{}
		if err := c.Get(ctx, types.NamespacedName{Namespace: namespace, Name: oci.ImagePullSecret.Name}, sec); err != nil {
			log.Error(err, "failed to retrieve authentication information")
			// todo report as a condition on the addon
			return auth, err
		}
		// read the JSON from the secret
		var authConfig DockerConfigJSON
		if err := json.Unmarshal(sec.Data[".dockerconfigjson"], &authConfig); err != nil {
			log.Error(err, "failed to unmarshall data from .dockerconfigjson")
			return auth, err
		}
		for k, v := range authConfig.Auths {
			log.V(1).Info("adding authentication data from secret", "registry", k, "username", v.Username)
			auth = append(auth, remote.WithAuth(authn.FromConfig(v)))
		}
	}
	return auth, nil
}
