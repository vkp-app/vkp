package graph

import (
	"context"
	"errors"
	"github.com/go-logr/logr"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prommodel "github.com/prometheus/common/model"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/model"
	v1 "gitlab.dcas.dev/k8s/kube-glass/apiserver/pkg/cfg/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

const labelTenant = "paas.dcas.dev/tenant"

var (
	ErrUnauthorised   = errors.New("unauthorised")
	ErrForbidden      = errors.New("forbidden")
	ErrTenantNotReady = errors.New("tenant is not ready")
)

type KubeOpts struct {
	GroupPrefix    string
	UsernamePrefix string
}

type Resolver struct {
	client.Client
	Scheme *runtime.Scheme

	prometheusAPI    promv1.API
	prometheusConfig *v1.PrometheusConfig

	dexURL string
	dexCA  string

	kubeOpts KubeOpts
}

func NewResolver(ctx context.Context, client client.Client, scheme *runtime.Scheme, prometheus promv1.API, prometheusConfig *v1.PrometheusConfig, kubeOpts KubeOpts, dexURL string, dexCA string) (*Resolver, error) {
	log := logr.FromContextOrDiscard(ctx)
	var caData string
	if dexCA != "" {
		log.Info("reading dex CA file", "path", dexCA)
		data, err := os.ReadFile(dexCA)
		if err != nil {
			log.Error(err, "failed to read dex CA file", "path", dexCA)
			return nil, err
		}
		caData = string(data)
	}

	return &Resolver{
		Client:           client,
		Scheme:           scheme,
		prometheusAPI:    prometheus,
		prometheusConfig: prometheusConfig,
		dexURL:           dexURL,
		dexCA:            caData,
		kubeOpts:         kubeOpts,
	}, nil
}

func (r *Resolver) GetMetrics(ctx context.Context, tenant, cluster string) ([]model.Metric, error) {
	srp := strings.NewReplacer(
		"{namespace}", tenant,
		"{cluster}", cluster,
	)
	var metrics []model.Metric
	for _, m := range r.prometheusConfig.ClusterMetrics {
		metric := srp.Replace(m.Metric)
		// fetch metrics from prometheus. If there's
		// an error, swallow it and return an empty list
		values, _ := r.GetMetric(ctx, metric)
		if values == nil {
			values = []model.MetricValue{}
		}
		metrics = append(metrics, model.Metric{
			Name:   m.Name,
			Metric: m.Metric,
			Format: m.Format,
			Values: values,
		})
	}
	return metrics, nil
}

func (r *Resolver) GetMetric(ctx context.Context, promQL string) ([]model.MetricValue, error) {
	log := logr.FromContextOrDiscard(ctx)
	log.V(1).Info("preparing prometheus query", "promql", promQL)

	resp, _, err := r.prometheusAPI.QueryRange(ctx, promQL, promv1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	})
	if err != nil {
		log.Error(err, "failed to query prometheus")
		return nil, err
	}
	log.V(2).Info("received response from Prometheus", "Raw", resp)
	// cast the response into something we know
	data, ok := resp.(prommodel.Matrix)
	if !ok {
		log.Info("failed to cast response data into model.Matrix")
		return nil, errors.New("unexpected data type returned from Prometheus")
	}
	if len(data) == 0 {
		log.V(1).Info("received empty data from Prometheus")
		return []model.MetricValue{}, nil
	}
	// convert the response data
	// into our graphql format
	results := make([]model.MetricValue, len(data[0].Values))
	for i, d := range data[0].Values {
		results[i] = model.MetricValue{
			Time:  d.Timestamp.Unix(),
			Value: d.Value.String(),
		}
	}
	return results, err
}
