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

package idp

import (
	"context"
	"fmt"
	"github.com/dexidp/dex/api/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"

	idpv1 "gitlab.dcas.dev/k8s/kube-glass/operator/apis/idp/v1"
)

const (
	eventClientCreated      = "DexClientCreated"
	eventClientDeleted      = "DexClientDeleted"
	eventClientCreateFailed = "DexClientCreateFailed"
	eventClientUpdateFailed = "DexClientUpdateFailed"
	eventClientDeleteFailed = "DexClientDeleteFailed"
)

// OAuthClientReconciler reconciles a OAuthClient object
type OAuthClientReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Options  OAuthClientOptions
}

type OAuthClientOptions struct {
	DexGrpcAddr string
}

func DexClientId(or *idpv1.OAuthClient) string {
	// todo verify if duplicate client_id's will cause a collision
	return or.Spec.ClientID
}

//+kubebuilder:rbac:groups=idp.dcas.dev,resources=oauthclients,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=idp.dcas.dev,resources=oauthclients/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=idp.dcas.dev,resources=oauthclients/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *OAuthClientReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logging.FromContext(ctx)

	or := &idpv1.OAuthClient{}
	if err := r.Get(ctx, req.NamespacedName, or); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to retrieve client resource")
		return ctrl.Result{}, err
	}
	if or.DeletionTimestamp != nil {
		log.Info("skipping oauth client that has been marked for termination")
		if controllerutil.ContainsFinalizer(or, finalizer) {
			// terminate the resource
			if err := r.finalizeDexClient(ctx, or); err != nil {
				return ctrl.Result{}, err
			}

			// remove the finalizer
			controllerutil.RemoveFinalizer(or, finalizer)
			if err := r.Update(ctx, or); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// add our finalizer
	if !controllerutil.ContainsFinalizer(or, finalizer) {
		controllerutil.AddFinalizer(or, finalizer)
		if err := r.Update(ctx, or); err != nil {
			log.Error(err, "failed to add finalizer")
			return ctrl.Result{}, err
		}
	}

	if err := r.reconcileDexClient(ctx, or); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *OAuthClientReconciler) reconcileDexClient(ctx context.Context, or *idpv1.OAuthClient) error {
	log := logging.FromContext(ctx)
	log.Info("reconciling OAuthClient", "addr", r.Options.DexGrpcAddr)

	// fetch the secret
	sc := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: or.GetNamespace(), Name: or.Spec.ClientSecretRef.Name}, sc); err != nil {
		log.Error(err, "failed to fetch dex secret")
		return err
	}

	// establish a connection to the Dex API
	conn, err := grpc.DialContext(ctx, r.Options.DexGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err, "failed to establish gRPC connection to Dex")
		return err
	}
	defer conn.Close()
	dexClient := api.NewDexClient(conn)
	oauthClient := &api.Client{
		Id:           DexClientId(or),
		Name:         fmt.Sprintf("%s/%s", or.GetNamespace(), or.GetName()),
		Secret:       string(sc.Data[or.Spec.ClientSecretRef.Key]),
		RedirectUris: or.Spec.RedirectURIs,
		TrustedPeers: or.Spec.TrustedPeers,
		Public:       or.Spec.Public,
		LogoUrl:      or.Spec.LogoURL,
	}
	// create the client
	resp, err := dexClient.CreateClient(ctx, &api.CreateClientReq{Client: oauthClient})
	if err != nil {
		r.Recorder.Eventf(or, corev1.EventTypeWarning, eventClientCreateFailed, `Failed to create Dex client "%s": %s`, DexClientId(or), err)
		log.Error(err, "failed to create Dex client")
		return err
	}
	if !resp.AlreadyExists {
		r.Recorder.Eventf(or, corev1.EventTypeNormal, eventClientCreated, `Successfully created Dex client "%s"`, DexClientId(or))
		return nil
	}
	// reconcile the client.
	// we can't do a diff here because we have no easy way
	// of fetching the current client data from Dex
	log.Info("patching Dex client")
	_, err = dexClient.UpdateClient(ctx, &api.UpdateClientReq{
		Id:           oauthClient.Id,
		RedirectUris: oauthClient.RedirectUris,
		TrustedPeers: oauthClient.TrustedPeers,
		Name:         oauthClient.Name,
		LogoUrl:      oauthClient.LogoUrl,
	})
	if err != nil {
		r.Recorder.Eventf(or, corev1.EventTypeWarning, eventClientUpdateFailed, `Failed to update Dex client "%s": %s`, DexClientId(or), err)
		log.Error(err, "failed to update Dex client")
		return err
	}
	return nil
}

func (r *OAuthClientReconciler) finalizeDexClient(ctx context.Context, or *idpv1.OAuthClient) error {
	log := logging.FromContext(ctx)
	log.Info("finalizing dex client", "addr", r.Options.DexGrpcAddr)
	// establish a connection to the Dex API
	conn, err := grpc.DialContext(ctx, r.Options.DexGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(err, "failed to establish gRPC connection to Dex")
		return err
	}
	defer conn.Close()
	dexClient := api.NewDexClient(conn)
	_, err = dexClient.DeleteClient(ctx, &api.DeleteClientReq{
		Id: DexClientId(or),
	})
	if err != nil {
		r.Recorder.Eventf(or, corev1.EventTypeWarning, eventClientDeleteFailed, `Failed to delete Dex client "%s": %s`, DexClientId(or), err)
		log.Error(err, "failed to delete Dex client")
		return err
	}
	r.Recorder.Eventf(or, corev1.EventTypeNormal, eventClientDeleted, `Successfully deleted Dex client "%s"`, DexClientId(or))
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OAuthClientReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&idpv1.OAuthClient{}).
		Complete(r)
}
