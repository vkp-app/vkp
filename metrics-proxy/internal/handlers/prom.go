package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/vkp-app/vkp/metrics-proxy/internal/promutil"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type contextKeyPromLabels int

const (
	labelNamespace contextKeyPromLabels = iota
	labelPods      contextKeyPromLabels = iota
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
		query = strings.ReplaceAll(query, fmt.Sprintf(`"%s"`, namespace), fmt.Sprintf(`"%s"`, tenant))
		query = strings.ReplaceAll(query, `pod="`, `pod=~"`)
		for i := range pods {
			query = strings.ReplaceAll(query, pods[i], fmt.Sprintf("%s-.*", pods[i]))
		}
		r.URL.RawQuery = "query=" + url.QueryEscape(query)
		log.Info("completed query rewrite", "query", query)

		ctx = context.WithValue(ctx, labelNamespace, namespace)
		ctx = context.WithValue(ctx, labelPods, pods)

		// continue as normal
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PrometheusRewriteResponse(r *http.Response) error {
	ctx := r.Request.Context()
	log := logr.FromContextOrDiscard(ctx)
	// read the response and decode it
	// if it has been gzipped
	var wasGzip bool
	var reader io.ReadCloser
	var err error
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		log.V(5).Info("decompressing response as it was gzipped")
		reader, err = gzip.NewReader(r.Body)
		wasGzip = true
	default:
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		log.Error(err, "failed to read response body")
		return err
	}
	_ = reader.Close()

	log.V(6).Info("dumping response body", "raw", string(body), "headers", r.Header)

	data := string(body)
	pods := ctx.Value(labelPods).([]string)
	log.V(3).Info("reversing pod labels", "pods", pods)
	data = reversePods(data, pods)

	log.V(6).Info("dumping new response body", "raw", data)

	// create the new body that we can respond with
	nb := io.NopCloser(strings.NewReader(data))
	if wasGzip {
		log.V(5).Info("compressing new response to match original")
		buf := bytes.NewBuffer(nil)
		w := gzip.NewWriter(buf)
		_, _ = io.Copy(w, nb)
		_ = w.Close()
		nb = io.NopCloser(buf)
	}
	r.Body = nb
	r.ContentLength = int64(len(body))
	r.Header.Set("Content-Length", strconv.Itoa(int(r.ContentLength)))

	return nil
}

var regexpPod = regexp.MustCompile(`pod=.?"([^"]+)"`)
var regexpNamespace = regexp.MustCompile(`namespace=.?"([^"]+)"`)
var errNoMatch = errors.New("could not locate label in query")

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

func reversePods(data string, pods []string) string {
	for _, p := range pods {
		rp := regexp.MustCompile(fmt.Sprintf(`%s([^"]*)`, p))
		data = rp.ReplaceAllString(data, p)
	}
	return data
}
