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
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"gitlab.com/autokubeops/serverless"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/generated"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/userctx"
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

type environment struct {
	Port     int `envconfig:"PORT" default:"8080"`
	LogLevel int `split_words:"true"`

	Otel struct {
		Enabled    bool    `split_words:"true"`
		SampleRate float64 `split_words:"true"`
	}

	Metrics struct {
		PrometheusAddr string `split_words:"true" required:"true"`
	}
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(paasv1alpha1.AddToScheme(scheme))
}

func main() {
	var e environment
	envconfig.MustProcess("api", &e)

	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(e.LogLevel * -1))

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
		Enabled:     e.Otel.Enabled,
		ServiceName: "glass-apiserver",
		SampleRate:  e.Otel.SampleRate,
	})
	if err != nil {
		log.Error(err, "failed to configure OpenTelemetry")
		os.Exit(1)
		return
	}

	promClient, err := api.NewClient(api.Config{Address: e.Metrics.PrometheusAddr})
	if err != nil {
		log.Error(err, "failed to create Prometheus client")
		os.Exit(1)
		return
	}

	// configure graphql
	resolver := &graph.Resolver{
		Client:     kubeClient,
		Scheme:     scheme,
		Prometheus: promv1.NewAPI(promClient),
	}
	c := generated.Config{Resolvers: resolver}
	c.Directives.HasUser = graph.HasUser
	c.Directives.HasAdmin = resolver.HasAdmin
	srv := handler.New(generated.NewExecutableSchema(c))
	srv.AddTransport(transport.POST{})

	// configure routing
	router := mux.NewRouter()
	router.Use(otel.Middleware(), logging.Middleware(log), metrics.Middleware(), userctx.Middleware())
	router.Handle("/metrics", prom)
	router.Handle("/api/v1/graphql", playground.Handler("GraphQL Playground", "/api/v1/query"))
	router.Handle("/api/v1/query", srv)

	// start the server
	serverless.NewBuilder(router).
		WithPort(e.Port).
		WithLogger(log).
		Run()
}
