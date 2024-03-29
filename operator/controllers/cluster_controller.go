/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	certv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	pgov1beta1 "github.com/crunchydata/postgres-operator/pkg/apis/postgres-operator.crunchydata.com/v1beta1"
	vclusterv1alpha1 "github.com/loft-sh/cluster-api-provider-vcluster/api/v1alpha1"
	idpv1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/idp/v1"
	"gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/cluster"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers/tenant"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"reflect"
	capiv1betav1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

const (
	eventDatabaseUserCreated      = "DatabaseUserCreated"
	eventDatabaseUserCreateFailed = "DatabaseUserCreateFailed"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	DexCA    string
	Options  ClusterOptions
}

//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=clusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=paas.dcas.dev,resources=appliedclusterversions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=postgres-operator.crunchydata.com,resources=postgresclusters,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cert-manager.io,resources=issuers,verbs=get;list;watch
//+kubebuilder:rbac:groups=cert-manager.io,resources=clusterissuers,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx).WithValues("cluster", req.NamespacedName)
	log.Info("reconciling Cluster")

	cr := &v1alpha1.Cluster{}
	if err := r.Get(ctx, req.NamespacedName, cr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve cluster resource")
		return ctrl.Result{}, err
	}
	if cr.DeletionTimestamp != nil {
		log.Info("skipping cluster that has been marked for deletion")
		if controllerutil.ContainsFinalizer(cr, finalizer) {
			// remove the finalizer since we were
			// able to successfully delete external
			// resources
			controllerutil.RemoveFinalizer(cr, finalizer)
			if err := r.Update(ctx, cr); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// backwards compat hack for #76
	if cr.Annotations[v1alpha1.LabelClusterID] == "" || cr.Annotations[v1alpha1.LabelClusterDomain] == "" {
		if cr.Annotations == nil {
			cr.Annotations = map[string]string{}
		}
		cr.Annotations[v1alpha1.LabelClusterID] = cr.Status.ClusterID
		cr.Annotations[v1alpha1.LabelClusterDomain] = cr.Status.ClusterDomain
		if err := r.Update(ctx, cr); err != nil {
			log.Error(err, "failed to update annotations")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil

	}
	cr.Status.WebURL = fmt.Sprintf("console.%s.%s", cr.Status.ClusterID, cr.Status.ClusterDomain)
	cr.Status.ClusterID = cr.Annotations[v1alpha1.LabelClusterID]
	cr.Status.ClusterDomain = cr.Annotations[v1alpha1.LabelClusterDomain]

	if err := r.safeguardHAClusters(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	res, conn, err := r.reconcileDatabase(ctx, cr)
	if err != nil || res.Requeue {
		return res, err
	}
	if err := r.reconcileDexSecret(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileAppliedClusterVersion(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileVCluster(ctx, cr, conn); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileCertificate(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileCluster(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileIngress(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.reconcileDexClient(ctx, cr); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Status().Update(ctx, cr); err != nil {
		log.Error(err, "failed to update status")
		return ctrl.Result{}, err
	}

	// add our finalizer
	if !controllerutil.ContainsFinalizer(cr, finalizer) {
		controllerutil.AddFinalizer(cr, finalizer)
		if err := r.Update(ctx, cr); err != nil {
			log.Error(err, "failed to add finalizer")
			return ctrl.Result{}, err
		}
	}

	// always requeue after 1 hour so that
	// we can perform upgrades during maintenance
	// windows.
	return ctrl.Result{RequeueAfter: time.Hour}, nil
}

func (r *ClusterReconciler) safeguardHAClusters(ctx context.Context, cr *v1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.V(1).Info("validating cluster HA mode")
	if !cr.Spec.HA.Enabled {
		return nil
	}
	// make sure we don't reconcile a cluster if
	// it's in an invalid state
	if cr.Spec.HA.Enabled && !r.Options.AllowHA {
		return fmt.Errorf("unable to process HA cluster when HA-support is not enabled")
	}
	return nil
}

func (r *ClusterReconciler) reconcileVCluster(ctx context.Context, cr *v1alpha1.Cluster, haConnectionString string) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling vcluster")

	latestVersion, err := r.getClusterTrack(ctx, cr)
	if err != nil {
		return err
	}

	vcluster, err := cluster.VCluster(ctx, cr, latestVersion, r.DexCA != "", tenant.CASecretName(cr.GetNamespace()), haConnectionString)
	if err != nil {
		return err
	}

	found := &vclusterv1alpha1.VCluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cr.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := ctrl.SetControllerReference(cr, vcluster, r.Scheme); err != nil {
				log.Error(err, "failed to set controller reference")
				return err
			}
			if err := r.Create(ctx, vcluster); err != nil {
				log.Error(err, "failed to create vcluster")
				return err
			}
			return nil
		}
		return err
	}
	cr.Status.Phase = vcluster.Status.Phase
	cr.Status.KubeURL = vcluster.Spec.ControlPlaneEndpoint.Host
	cr.Status.KubeVersion = latestVersion.Status.VersionNumber.String()
	cr.Status.PlatformVersion = latestVersion.GetName()
	// reconcile by forcibly overwriting
	// any changes
	if !reflect.DeepEqual(vcluster.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, vcluster)
	}
	return nil
}

func (r *ClusterReconciler) getClusterTrack(ctx context.Context, cr *v1alpha1.Cluster) (*v1alpha1.ClusterVersion, error) {
	log := logging.FromContext(ctx)
	log.V(1).Info("fetching most-appropriate cluster version")

	acv := &v1alpha1.AppliedClusterVersion{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cr.GetName()}, acv); err != nil {
		log.Error(err, "failed to retrieve AppliedClusterVersion resource")
		return nil, err
	}

	track := &v1alpha1.ClusterVersion{}
	if err := r.Get(ctx, types.NamespacedName{Name: acv.Status.VersionRef.Name}, track); err != nil {
		log.Error(err, "failed to retrieve ClusterVersion")
		return nil, err
	}
	return track, nil
}

func (r *ClusterReconciler) reconcileCluster(ctx context.Context, cr *v1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling cluster")

	capiCluster := cluster.Cluster(cr)

	found := &capiv1betav1.Cluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cr.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := ctrl.SetControllerReference(cr, capiCluster, r.Scheme); err != nil {
				log.Error(err, "failed to set controller reference")
				return err
			}
			if err := r.Create(ctx, capiCluster); err != nil {
				log.Error(err, "failed to create cluster")
				return err
			}
			return nil
		}
		return err
	}
	// reconcile by forcibly overwriting
	// any changes
	if !reflect.DeepEqual(capiCluster.Spec, found.Spec) {
		if err := ctrl.SetControllerReference(cr, capiCluster, r.Scheme); err != nil {
			log.Error(err, "failed to set controller reference")
			return err
		}
		return r.SafeUpdate(ctx, found, capiCluster)
	}
	return nil
}

func (r *ClusterReconciler) reconcileCertificate(ctx context.Context, cr *v1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling Certificate")

	cert := cluster.Certificate(cr, r.Options.RootCAIssuerName, r.Options.RootCAIssuerKind)

	found := &certv1.Certificate{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cert.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(cr, cert, r.Scheme)
			if err := r.Create(ctx, cert); err != nil {
				log.Error(err, "failed to create Certificate")
				return err
			}
			return nil
		}
		return err
	}
	// reconcile by forcibly overwriting
	// any changes
	if !reflect.DeepEqual(cert.Spec, found.Spec) {
		_ = ctrl.SetControllerReference(cr, cert, r.Scheme)
		return r.SafeUpdate(ctx, found, cert)
	}
	return nil
}

func (r *ClusterReconciler) reconcileAppliedClusterVersion(ctx context.Context, cr *v1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling AppliedClusterVersion")

	acv := cluster.AppliedClusterVersion(cr)

	found := &v1alpha1.AppliedClusterVersion{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: acv.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(cr, acv, r.Scheme)
			if err := r.Create(ctx, acv); err != nil {
				log.Error(err, "failed to create AppliedClusterVersion")
				return err
			}
			return nil
		}
		return err
	}
	// don't patch anything since we want it to
	// be user-configurable
	return nil
}

func (r *ClusterReconciler) reconcileIngress(ctx context.Context, cr *v1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling ingress")

	ing := cluster.Ingress(cr)

	found := &netv1.Ingress{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cr.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(cr, ing, r.Scheme)
			if err := r.Create(ctx, ing); err != nil {
				log.Error(err, "failed to create ingress")
				return err
			}
			return nil
		}
		return err
	}
	// reconcile by forcibly overwriting
	// any changes
	if !reflect.DeepEqual(ing.Spec, found.Spec) {
		_ = ctrl.SetControllerReference(cr, ing, r.Scheme)
		return r.SafeUpdate(ctx, found, ing)
	}
	return nil
}

func (r *ClusterReconciler) reconcileDexSecret(ctx context.Context, cr *v1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling dex secret")

	sec := cluster.DexSecret(cr, r.DexCA)

	found := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: sec.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			if err := ctrl.SetControllerReference(cr, sec, r.Scheme); err != nil {
				log.Error(err, "failed to set controller reference")
				return err
			}
			if err := r.Create(ctx, sec); err != nil {
				log.Error(err, "failed to create dex secret")
				return err
			}
			return nil
		}
		return err
	}
	_ = ctrl.SetControllerReference(cr, sec, r.Scheme)
	if found.Data == nil {
		log.Info("skipping dex secret as .data is nil")
		found.Data = sec.Data
		return r.Update(ctx, found)
	}
	// validate that the client id exists
	if err := r.updateSecretData(ctx, cluster.DexKeyID, found, sec); err != nil {
		return err
	}
	// validate that the client secret exists
	if err := r.updateSecretData(ctx, cluster.DexKeySecret, found, sec); err != nil {
		return err
	}
	// validate and update the CA certificate
	if err := r.updateSecretData(ctx, cluster.DexKeyCA, found, sec); err != nil {
		return err
	}
	return nil
}

func (r *ClusterReconciler) updateSecretData(ctx context.Context, key string, found, sec *corev1.Secret) error {
	log := logging.FromContext(ctx)
	existing := sec.Data[key]
	if val, ok := found.Data[key]; !ok || val == nil {
		found.Data[key] = existing
		if err := r.Update(ctx, found); err != nil {
			log.Error(err, "failed to update dex secret")
			return err
		}
	}
	return nil
}

func (r *ClusterReconciler) reconcileDatabase(ctx context.Context, cr *v1alpha1.Cluster) (ctrl.Result, string, error) {
	log := logging.FromContext(ctx)
	log.Info("reconciling database")

	if !cr.Spec.HA.Enabled {
		log.V(1).Info("skipping non-HA cluster")
		return ctrl.Result{}, "", nil
	}

	if !r.Options.AllowHA {
		log.V(1).Info("skipping HA cluster due to policy")
		return ctrl.Result{}, "", nil
	}

	db := &pgov1beta1.PostgresCluster{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: r.Options.PostgresResourceNamespace, Name: r.Options.PostgresResourceName}, db); err != nil {
		log.Error(err, "failed to fetch HA database resource")
		return ctrl.Result{}, "", err
	}

	if r.Options.UseHANonce && cr.Status.ClusterID == "" {
		log.V(1).Info("waiting before creating database since the cluster ID has not been assigned")
		return ctrl.Result{Requeue: true}, "", nil
	}
	username := fmt.Sprintf("%s-%s", cr.GetNamespace(), cr.GetName())
	database := cr.Status.ClusterDatabase
	if database == "" {
		// make sure to add a nonce so that recreating the cluster
		// with the same name doesn't cause db/encryption issues
		if r.Options.UseHANonce {
			log.V(2).Info("attaching cluster ID to database")
			database = fmt.Sprintf("%s-%s", username, cr.Status.ClusterID)
		} else {
			database = username
		}
		// update the status
		cr.Status.ClusterDatabase = database
		if err := r.Status().Update(ctx, cr); err != nil {
			log.Error(err, "failed to update Cluster database status")
			return ctrl.Result{}, "", err
		}
		return ctrl.Result{Requeue: true}, "", nil
	}
	// check if our user has been created.
	for _, u := range db.Spec.Users {
		// if our user has been added, we're
		// done here
		if string(u.Name) == username {
			log.Info("successfully located our user within the database resource", "username", username)
			return r.getConnectionString(ctx, username)
		}
	}
	// otherwise, we need to add our user
	log.Info("adding our user to the database resource", "username", username, "database", database)
	db.Spec.Users = append(db.Spec.Users, pgov1beta1.PostgresUserSpec{
		Name:      pgov1beta1.PostgresIdentifier(username),
		Databases: []pgov1beta1.PostgresIdentifier{pgov1beta1.PostgresIdentifier(database)},
		Password: &pgov1beta1.PostgresPasswordSpec{
			Type: pgov1beta1.PostgresPasswordTypeAlphaNumeric,
		},
	})
	if err := r.Update(ctx, db); err != nil {
		r.Recorder.Eventf(cr, corev1.EventTypeWarning, eventDatabaseUserCreateFailed, `Failed to create database and user "%s": %s`, database, err.Error())
		log.Error(err, "failed to update HA database resource")
		return ctrl.Result{}, "", err
	}
	r.Recorder.Eventf(cr, corev1.EventTypeNormal, eventDatabaseUserCreated, `Successfully created database and user "%s"`, database)

	// fetch the connection string
	return r.getConnectionString(ctx, username)
}

func (r *ClusterReconciler) getConnectionString(ctx context.Context, username string) (ctrl.Result, string, error) {
	log := logging.FromContext(ctx)
	log.Info("fetching connection string")
	// fetch the connection string
	sec := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: r.Options.PostgresResourceNamespace, Name: fmt.Sprintf("%s-pguser-%s", r.Options.PostgresResourceName, username)}, sec); err != nil {
		// if the secret doesn't exist yet,
		// wait until it does
		if errors.IsNotFound(err) {
			log.Info("pguser secret does not exist yet")
			return ctrl.Result{Requeue: true}, "", nil
		}
		log.Error(err, "failed to retrieve pguser secret")
		return ctrl.Result{}, "", err
	}

	conn := string(sec.Data[secretKeyDbConn])
	if conn == "" {
		log.Info("pguser connection string is not ready yet")
		return ctrl.Result{Requeue: true}, "", nil
	}
	return ctrl.Result{}, conn, nil
}

func (r *ClusterReconciler) reconcileDexClient(ctx context.Context, cr *v1alpha1.Cluster) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling dex client")

	// fetch the secret
	sc := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: cluster.DexSecretName(cr.GetName())}, sc); err != nil {
		log.Error(err, "failed to fetch dex secret")
		return err
	}

	dexClient := cluster.DexOAuthClient(cr, sc)

	found := &idpv1.OAuthClient{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: cr.GetNamespace(), Name: dexClient.GetName()}, found); err != nil {
		if errors.IsNotFound(err) {
			_ = ctrl.SetControllerReference(cr, dexClient, r.Scheme)
			if err := r.Create(ctx, dexClient); err != nil {
				log.Error(err, "failed to create dex client")
				return err
			}
			return nil
		}
		return err
	}
	_ = ctrl.SetControllerReference(cr, dexClient, r.Scheme)
	if !reflect.DeepEqual(dexClient.Spec, found.Spec) {
		return r.SafeUpdate(ctx, found, dexClient)
	}
	return nil
}

// SafeUpdate calls Update with hacks required to ensure that
// the update is applied correctly.
//
// https://github.com/argoproj/argo-cd/issues/3657
func (r *ClusterReconciler) SafeUpdate(ctx context.Context, old, new client.Object, option ...client.UpdateOption) error {
	new.SetResourceVersion(old.GetResourceVersion())
	return r.Update(ctx, new, option...)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Cluster{}).
		Owns(&v1alpha1.AppliedClusterVersion{}).
		Owns(&vclusterv1alpha1.VCluster{}).
		Owns(&capiv1betav1.Cluster{}).
		Owns(&netv1.Ingress{}).
		Owns(&corev1.Secret{}).
		Owns(&idpv1.OAuthClient{}).
		Owns(&certv1.Certificate{}).
		Complete(r)
}
