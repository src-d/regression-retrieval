package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/src-d/regression-core"
	"gopkg.in/src-d/go-log.v1"
)

const (
	WSeconds  = "regression_retrieval_w_avg_seconds"
	SSeconds  = "regression_retrieval_s_avg_seconds"
	USeconds  = "regression_retrieval_u_avg_seconds"
	MemoryMiB = "regression_retrieval_mem_avg_mib"
)

// labels represents labels to be linked on observations
var labels = []string{"util", "version", "name", "branch", "commit"}

// metrics is a map of name to type SummaryVec pointer
type metrics map[string]*prometheus.SummaryVec

// PromClient is the wrapper around pusher that also keeps metrics
type PromClient struct {
	pusher  *push.Pusher
	metrics metrics
	util    string
}

// TODO(lwsanty): one day it possibly could be a part of regression-core
// NewPromClient inits new pusher, creates metrics and adds them to the collector
func NewPromClient(util string, p regression.PromConfig) *PromClient {
	pusher := push.New(p.Address, p.Job)
	log.Debugf("adding metrics to the pusher")

	metrics := getMetrics(labels)
	for _, m := range metrics {
		pusher.Collector(m)
	}
	return &PromClient{
		pusher:  pusher,
		metrics: metrics,
		util:    util,
	}
}

// Dump does observations and adds metrics to the pusher
func (p *PromClient) Dump(res *regression.Result, version, name, branch, commit string) error {
	labelValues := []string{p.util, version, name, branch, commit}
	observe := func(metric string, value float64) {
		p.metrics[metric].WithLabelValues(labelValues...).Observe(value)
	}
	observe(WSeconds, res.Wtime.Seconds())
	observe(SSeconds, res.Stime.Seconds())
	observe(USeconds, res.Utime.Seconds())
	observe(MemoryMiB, toMiB(res.Memory))

	log.Debugf("pushing metrics")
	return p.pusher.Add()
}

func getMetrics(labels []string) metrics {
	return metrics{
		WSeconds:  getMetric(WSeconds, labels),
		SSeconds:  getMetric(SSeconds, labels),
		USeconds:  getMetric(USeconds, labels),
		MemoryMiB: getMetric(MemoryMiB, labels),
	}
}

func getMetric(name string, labels []string) *prometheus.SummaryVec {
	return prometheus.NewSummaryVec(
		prometheus.SummaryOpts{Name: name},
		labels,
	)
}

func toMiB(i int64) float64 {
	return float64(i) / float64(1024*1024)
}
