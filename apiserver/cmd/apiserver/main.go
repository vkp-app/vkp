package main

import (
	"context"
	"github.com/djcass44/go-utils/logging"
	"github.com/djcass44/go-utils/otel"
	"github.com/djcass44/go-utils/otel/metrics"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/autokubeops/serverless"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type environment struct {
	Port     int `envconfig:"PORT" default:"8080"`
	LogLevel int `split_words:"true"`

	Otel struct {
		Enabled    bool    `split_words:"true"`
		SampleRate float64 `split_words:"true"`
	}
}

func main() {
	var e environment
	envconfig.MustProcess("api", &e)

	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(e.LogLevel * -1))

	log, ctx := logging.NewZap(context.TODO(), zc)

	// configure metrics and tracing
	prom := metrics.MustNewDefault(ctx)
	err := otel.Build(ctx, otel.Options{
		Enabled:     e.Otel.Enabled,
		ServiceName: "glass-apiserver",
		SampleRate:  e.Otel.SampleRate,
	})
	if err != nil {
		log.Error(err, "failed to configure OpenTelemetry")
		os.Exit(1)
		return
	}

	// configure routing
	router := mux.NewRouter()
	router.Use(otel.Middleware(), logging.Middleware(log), metrics.Middleware())
	router.Handle("/metrics", prom)

	// start the server
	serverless.NewBuilder(router).
		WithPort(e.Port).
		WithLogger(log).
		Run()
}
