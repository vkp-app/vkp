package controllers

import (
	"context"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/platform"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *PlatformReconciler) reconcileApiDeployment(ctx context.Context, pr *paasv1alpha1.Platform) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling api server deployment")

	dp := platform.ApiDeployment(pr)

	found := &appsv1.Deployment{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: dp.GetNamespace(), Name: dp.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, dp); err != nil {
				log.Error(err, "failed to create api server deployment")
				return err
			}
			return nil
		}
		return err
	}

	// reconcile changes
	if !reflect.DeepEqual(dp.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, dp)
	}
	return nil
}

func (r *PlatformReconciler) reconcileApiService(ctx context.Context, pr *paasv1alpha1.Platform) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling api server service")

	svc := platform.ApiService(pr)

	found := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: svc.GetNamespace(), Name: svc.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, svc); err != nil {
				log.Error(err, "failed to create api server service")
				return err
			}
			return nil
		}
		return err
	}

	// reconcile changes
	if !reflect.DeepEqual(svc.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, svc)
	}
	return nil
}

func (r *PlatformReconciler) reconcileApiConfig(ctx context.Context, pr *paasv1alpha1.Platform) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling api server config")

	cm := platform.ApiConfig(pr)

	found := &corev1.ConfigMap{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cm.GetNamespace(), Name: cm.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, cm); err != nil {
				log.Error(err, "failed to create api server config")
				return err
			}
			return nil
		}
		return err
	}

	// reconcile changes
	if !reflect.DeepEqual(cm.Data, found.Data) {
		return r.SafeUpdate(ctx, found, cm)
	}
	return nil
}

func (r *PlatformReconciler) reconcileApiIngress(ctx context.Context, pr *paasv1alpha1.Platform) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling api server ingress")

	ing := platform.ApiIngress(pr)

	found := &netv1.Ingress{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: ing.GetNamespace(), Name: ing.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, ing); err != nil {
				log.Error(err, "failed to create api server ingress")
				return err
			}
			return nil
		}
		return err
	}

	// reconcile changes
	if !reflect.DeepEqual(ing.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, ing)
	}
	return nil
}
