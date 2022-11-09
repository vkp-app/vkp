package graph

import (
	"context"
	"errors"
	"github.com/go-logr/logr"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prommodel "github.com/prometheus/common/model"
	"gitlab.dcas.dev/k8s/kube-glass/apiserver/internal/graph/model"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const labelTenant = "paas.dcas.dev/tenant"

var ErrUnauthorised = errors.New("unauthorised")
var ErrTenantNotReady = errors.New("tenant is not ready")

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client.Client
	Scheme     *runtime.Scheme
	Prometheus promv1.API
}

func (r *Resolver) GetMetric(ctx context.Context, promQL string) ([]model.MetricValue, error) {
	log := logr.FromContextOrDiscard(ctx)
	log.V(1).Info("preparing prometheus query", "promql", promQL)

	resp, _, err := r.Prometheus.QueryRange(ctx, promQL, promv1.Range{
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
