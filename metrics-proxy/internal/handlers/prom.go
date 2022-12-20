package handlers

import (
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/vkp-app/vkp/metrics-proxy/internal/promutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func PrometheusRewrite(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cluster := ctx.Value(RequestingCluster).(string)
		tenant := ctx.Value(RequestingTenant).(string)

		log := logr.FromContextOrDiscard(ctx).WithValues("tenant", tenant, "cluster", cluster)

		query := r.URL.Query().Get("query")
		log.Info("processing query", "query", query)

		// figure out what type of query we're running
		namespace, pods, err := getQueryInfo(query)
		if err != nil {
			log.Error(err, "failed to parse query")
			promutil.RespondErr(ctx, w, err, http.StatusBadRequest)
			return
		}

		// rewrite the query
		query = strings.ReplaceAll(query, namespace, tenant)
		for i := range pods {
			query = strings.ReplaceAll(query, pods[i], fmt.Sprintf("%s-.*", pods[i]))
		}
		r.URL.RawQuery = "query=" + url.QueryEscape(query)
		log.Info("completed query rewrite", "query", query)

		// continue as normal
		h.ServeHTTP(w, r)
	})
}

var regexpPod = regexp.MustCompile(`pod=.?"([^"]+)"`)
var regexpNamespace = regexp.MustCompile(`namespace=.?"([^"]+)"`)
var errNoMatch = errors.New("could not locate label in query")

const (
	queryContainerCPU    = "container_cpu_usage_seconds_total"
	queryContainerMemory = "container_memory_working_set_bytes"
)

func getQueryInfo(query string) (string, []string, error) {
	var ns string
	matches := regexpNamespace.FindStringSubmatch(query)
	if len(matches) < 2 {
		return "", nil, errNoMatch
	}
	ns = matches[1]
	matches = regexpPod.FindStringSubmatch(query)
	if len(matches) < 2 {
		return "", nil, errNoMatch
	}

	return ns, strings.Split(matches[1], "|"), nil
}
