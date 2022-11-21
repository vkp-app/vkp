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
	ctrl "sigs.k8s.io/controller-runtime"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *PlatformReconciler) reconcileDexDeployment(ctx context.Context, pr *paasv1alpha1.Platform) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling Dex server deployment")

	dp := platform.DexDeployment(pr)

	found := &appsv1.Deployment{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: dp.GetNamespace(), Name: dp.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(pr, dp, r.Scheme)
			if err := r.Create(ctx, dp); err != nil {
				log.Error(err, "failed to create Dex server deployment")
				return err
			}
			return nil
		}
		return err
	}
	_ = ctrl.SetControllerReference(pr, dp, r.Scheme)

	// reconcile changes
	if !reflect.DeepEqual(dp.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, dp)
	}
	return nil
}

func (r *PlatformReconciler) reconcileDexService(ctx context.Context, pr *paasv1alpha1.Platform) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling Dex server service")

	svc := platform.DexService(pr)

	found := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: svc.GetNamespace(), Name: svc.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(pr, svc, r.Scheme)
			if err := r.Create(ctx, svc); err != nil {
				log.Error(err, "failed to create Dex server service")
				return err
			}
			return nil
		}
		return err
	}
	_ = ctrl.SetControllerReference(pr, svc, r.Scheme)

	// reconcile changes
	if !reflect.DeepEqual(svc.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, svc)
	}
	return nil
}

func (r *PlatformReconciler) reconcileDexConfig(ctx context.Context, pr *paasv1alpha1.Platform) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling Dex server config")

	cm, err := platform.DexConfig(ctx, pr)
	if err != nil {
		return err
	}

	found := &corev1.ConfigMap{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cm.GetNamespace(), Name: cm.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(pr, cm, r.Scheme)
			if err := r.Create(ctx, cm); err != nil {
				log.Error(err, "failed to create Dex server config")
				return err
			}
			return nil
		}
		return err
	}
	_ = ctrl.SetControllerReference(pr, cm, r.Scheme)

	// reconcile changes
	if !reflect.DeepEqual(cm.Data, found.Data) {
		return r.SafeUpdate(ctx, found, cm)
	}
	return nil
}

func (r *PlatformReconciler) reconcileDexIngress(ctx context.Context, pr *paasv1alpha1.Platform) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling Dex server ingress")

	ing := platform.DexIngress(pr)

	found := &netv1.Ingress{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: ing.GetNamespace(), Name: ing.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(pr, ing, r.Scheme)
			if err := r.Create(ctx, ing); err != nil {
				log.Error(err, "failed to create Dex server ingress")
				return err
			}
			return nil
		}
		return err
	}

	_ = ctrl.SetControllerReference(pr, ing, r.Scheme)

	// reconcile changes
	if !reflect.DeepEqual(ing.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, ing)
	}
	return nil
}
