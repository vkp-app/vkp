package main

import (
	"context"
	"crypto/tls"
	"github.com/djcass44/go-utils/logging"
	"github.com/gorilla/mux"
	flag "github.com/spf13/pflag"
	"github.com/vkp-app/vkp/metrics-proxy/internal/handlers"
	"gitlab.com/autokubeops/serverless"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	// flags
	fLogLevel := flag.Int("v", 0, "level of logging information (higher is more).")
	fPort := flag.Int("port", 8080, "http port to listen on.")

	fPrometheusURL := flag.String("prometheus-url", "", "URL of the upstream Prometheus instance.")

	fCertFile := flag.String("tls-cert-file", "", "path to the TLS certificate file.")
	fKeyFile := flag.String("tls-key-file", "", "path to the TLS key file.")

	flag.Parse()

	// logging configuration
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(*fLogLevel * -1))

	log, _ := logging.NewZap(context.TODO(), zc)

	// validate the url
	uri, err := url.Parse(*fPrometheusURL)
	if err != nil {
		log.Error(err, "failed to parse prometheus url")
		os.Exit(1)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(uri)
	proxy.ModifyResponse = handlers.PrometheusRewriteResponse

	// configure routing
	router := mux.NewRouter()
	router.Use(logging.Middleware(log))
	router.Handle("/api/v1/query", handlers.TLS(handlers.PrometheusRewrite(proxy)))

	// start the server
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	serverless.NewBuilder(router).
		WithLogger(log).
		WithPort(*fPort).
		WithTLS(*fCertFile, *fKeyFile, tlsConfig).
		Run()
}
