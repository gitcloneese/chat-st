package redlock

import "xy3-proto/pkg/stat/metric"

const (
	redLockNamespace = "red_lock"
)

var (
	//nolint:promlinter
	_metricRedLockReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: redLockNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "red lock requests duration(ms).",
		Labels:    []string{"key"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
)
