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

package main

import (
	"context"
	"github.com/djcass44/go-utils/logging"
	"github.com/jnovack/flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	pgov1beta1 "github.com/crunchydata/postgres-operator/pkg/apis/postgres-operator.crunchydata.com/v1beta1"
	vclusterv1alpha1 "github.com/loft-sh/cluster-api-provider-vcluster/api/v1alpha1"
	capiv1betav1 "sigs.k8s.io/cluster-api/api/v1beta1"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	"gitlab.dcas.dev/k8s/kube-glass/operator/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(paasv1alpha1.AddToScheme(scheme))
	utilruntime.Must(vclusterv1alpha1.AddToScheme(scheme))
	utilruntime.Must(capiv1betav1.AddToScheme(scheme))
	utilruntime.Must(pgov1beta1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	// custom flags
	fDexCA := flag.String("dex-ca-file", "", "file that contains the Certificate Authority for Dex. Will fallback to the Kubernetes API CA if not set.")
	fLogLevel := flag.Int("v", 0, "log verbosity (higher is more).")
	fLogDebug := flag.Bool("log-debug", false, "puts logging in development mode.")

	// cluster controller configuration
	clusterOpts := controllers.ClusterOptions{}
	flag.StringVar(&clusterOpts.DexGrpcAddr, "cluster-dex-grpc-addr", "", "grpc address of the Dex instance.")
	flag.BoolVar(&clusterOpts.AllowHA, "cluster-allow-ha", false, "determines whether HA control-planes should be allowed given the dependency on the Postgres Operator.")
	flag.StringVar(&clusterOpts.PostgresResourceName, "cluster-postgres-resource-name", "vkp", "name of the Postgres Operator Cluster resource to use for HA clusters.")
	flag.StringVar(&clusterOpts.PostgresResourceNamespace, "cluster-postgres-resource-namespace", os.Getenv("KUBERNETES_NAMESPACE"), "namespace of the Postgres Operator Cluster resource to use for HA clusters.")

	// tenant controller configuration
	tenantOpts := controllers.TenantOptions{}
	flag.BoolVar(&tenantOpts.SkipDefaultAddons, "tenant-skip-default-addons", false, "if enabled, will skip installation of cluster-wide addons.")
	flag.BoolVar(&tenantOpts.NamespaceOwnership, "tenant-namespace-ownership", true, "if enabled, Namespaces created will be owned by the Tenant resource.")
	fRootCA := flag.String("tenant-custom-ca-file", "", "file that contains one or more Certificate Authorities to be injected to all clusters.")

	flag.Parse()

	// configure logging
	var zc zap.Config
	if *fLogDebug {
		zc = zap.NewDevelopmentConfig()
	} else {
		zc = zap.NewProductionConfig()
	}
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(*fLogLevel * -1))

	log, _ := logging.NewZap(context.Background(), zc)

	// read the Dex CA
	var dexCA string
	if *fDexCA != "" {
		log.Info("reading dex CA", "path", *fDexCA)
		data, err := os.ReadFile(*fDexCA)
		if err != nil {
			log.Error(err, "failed to read dex CA file", "path", *fDexCA)
			os.Exit(1)
			return
		}
		dexCA = string(data)
	}
	if *fRootCA != "" {
		log.Info("reading root CA", "path", *fRootCA)
		data, err := os.ReadFile(*fRootCA)
		if err != nil {
			log.Error(err, "failed to read root CA file", "path", *fRootCA)
		}
		tenantOpts.CustomCAFile = string(data)
	}

	ctrl.SetLogger(log)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "aaac7635.dcas.dev",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.ClusterReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		DexCA:   dexCA,
		Options: clusterOpts,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Cluster")
		os.Exit(1)
	}
	if err = (&controllers.TenantReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Options: tenantOpts,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Tenant")
		os.Exit(1)
	}
	if err = (&controllers.ClusterAddonReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ClusterAddon")
		os.Exit(1)
	}
	if err = (&controllers.ClusterAddonBindingReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ClusterAddonBinding")
		os.Exit(1)
	}
	if err = (&controllers.ClusterVersionReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ClusterVersion")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
