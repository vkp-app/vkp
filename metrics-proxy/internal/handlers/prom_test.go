package handlers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrometheusRewrite_getQueryInfo(t *testing.T) {
	var cases = []struct {
		query     string
		namespace string
		pods      []string
		err       error
	}{
		{
			`sum by (pod,container) (
 container_memory_working_set_bytes{namespace="kube-system",pod=~"prometheus-adapter-7fffbc768b-zbgnm|coredns-57f58c9849-b9fl5",container!="",pod!=""}
)
`,
			"kube-system",
			[]string{"prometheus-adapter-7fffbc768b-zbgnm", "coredns-57f58c9849-b9fl5"},
			nil,
		},
		{
			`sum by (pod,container) (
 irate (
 container_cpu_usage_seconds_total{namespace="kube-system",pod=~"prometheus-adapter-7fffbc768b-zbgnm|coredns-57f58c9849-b9fl5",container!="",pod!=""}[4m]
 )
)
`,
			"kube-system",
			[]string{"prometheus-adapter-7fffbc768b-zbgnm", "coredns-57f58c9849-b9fl5"},
			nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.query, func(t *testing.T) {
			namespace, pods, err := getQueryInfo(tt.query)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.namespace, namespace)
			assert.ElementsMatch(t, tt.pods, pods)
		})
	}
}
