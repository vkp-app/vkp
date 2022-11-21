package controllers

import (
	"context"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/platform"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *PlatformReconciler) reconcileCommonSecret(ctx context.Context, pr *paasv1alpha1.Platform) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling common secret")

	sec := platform.CommonSecrets(pr)

	found := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: sec.GetNamespace(), Name: sec.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(pr, sec, r.Scheme)
			if err := r.Create(ctx, sec); err != nil {
				log.Error(err, "failed to create common secret")
				return err
			}
			return nil
		}
		return err
	}
	_ = ctrl.SetControllerReference(pr, sec, r.Scheme)

	// reconcile changes
	if found.Data == nil {
		found.Data = map[string][]byte{}
	}
	if val := found.Data[platform.SecretKeyOauthCookie]; val == nil {
		found.Data[platform.SecretKeyOauthCookie] = sec.Data[platform.SecretKeyOauthCookie]
		if err := r.Update(ctx, found); err != nil {
			return err
		}
	}
	if val := found.Data[platform.SecretKeyDexClientSecret]; val == nil {
		found.Data[platform.SecretKeyDexClientSecret] = sec.Data[platform.SecretKeyDexClientSecret]
		if err := r.Update(ctx, found); err != nil {
			return err
		}
	}
	return nil
}
