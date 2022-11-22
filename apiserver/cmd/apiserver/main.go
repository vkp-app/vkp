package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/djcass44/go-utils/logging"
	"github.com/djcass44/go-utils/otel"
	"github.com/djcass44/go-utils/otel/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	flag "github.com/spf13/pflag"
	"gitlab.com/autokubeops/serverless"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/generated"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/userctx"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/pkg/cfg"
	v1 "gitlab.dcas.dev/k8s/kube-glass/apiserver/pkg/cfg/v1"
	paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(paasv1alpha1.AddToScheme(scheme))
}

func main() {
	// flags
	fLogLevel := flag.Int("v", 0, "level of logging information (higher is more).")
	fPort := flag.Int("port", 8080, "http port to listen on.")

	fOtelEnabled := flag.Bool("otel-enabled", false, "enable exporting of OpenTelemetry traces.")
	fOtelServiceName := flag.String("otel-service-name", "glass-apiserver", "name to distinguish which service traces originate from.")
	fOtelSampleRate := flag.Float64("otel-sample-rate", 0, "percentage (0.0 - 1.0) of traces that should be exported.")

	fPrometheusURL := flag.String("prometheus-url", "", "URL of the management cluster's Prometheus server.")
	fPrometheusConfig := flag.String("prometheus-config-file", "", "File that contains the Prometheus configuration file.")
	fPrometheusMetrics := flag.Bool("prometheus-metrics", true, "Flag to indicate if Prometheus metrics should be exported.")

	fDexURL := flag.String("dex-url", "", "URL of the Dex instance.")
	fDexCA := flag.String("dex-ca-file", "", "File that contains the Certificate Authority for Dex. Will fallback to the Kubernetes API CA if not set.")

	flag.Parse()

	// logging configuration
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(*fLogLevel * -1))

	log, ctx := logging.NewZap(context.TODO(), zc)

	config := ctrl.GetConfigOrDie()
	kubeClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Error(err, "failed to build kube client")
		os.Exit(1)
		return
	}

	// configure metrics and tracing
	prom := metrics.MustNewDefault(ctx)
	err = otel.Build(ctx, otel.Options{
		Enabled:     *fOtelEnabled,
		ServiceName: *fOtelServiceName,
		SampleRate:  *fOtelSampleRate,
	})
	if err != nil {
		log.Error(err, "failed to configure OpenTelemetry")
		os.Exit(1)
		return
	}

	promClient, err := api.NewClient(api.Config{Address: os.ExpandEnv(*fPrometheusURL)})
	if err != nil {
		log.Error(err, "failed to create Prometheus client")
		os.Exit(1)
		return
	}

	// fetch configuration
	promConfig, err := cfg.Read[v1.PrometheusConfig](ctx, *fPrometheusConfig, v1.NewPrometheusConfig())
	if err != nil {
		log.Error(err, "failed to read Prometheus configuration file")
		os.Exit(1)
		return
	}

	// configure graphql
	resolver, err := graph.NewResolver(ctx, kubeClient, scheme, promv1.NewAPI(promClient), promConfig, *fDexURL, *fDexCA)
	if err != nil {
		log.Error(err, "failed to setup resolver")
		os.Exit(1)
		return
	}
	c := generated.Config{Resolvers: resolver}
	c.Directives.HasRole = resolver.HasRole
	c.Directives.HasClusterAccess = resolver.HasClusterAccess
	c.Directives.HasTenantAccess = resolver.HasTenantAccess
	srv := handler.New(generated.NewExecutableSchema(c))
	srv.AddTransport(transport.POST{})

	// configure routing
	router := mux.NewRouter()
	router.Use(otel.Middleware(), logging.Middleware(log), metrics.Middleware(), userctx.Middleware())
	if *fPrometheusMetrics {
		router.Handle("/metrics", prom)
	}
	router.Handle("/api/v1/graphql", playground.Handler("GraphQL Playground", "/api/v1/query"))
	router.Handle("/api/v1/query", srv)

	// start the server
	serverless.NewBuilder(router).
		WithPort(*fPort).
		WithLogger(log).
		Run()
}
