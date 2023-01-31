package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsClient struct {
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
}

// New inits the metrics and automatically registers the metrics with the default prometheus default register.
// When adding tags keep in mind tag cardinality, prometheus client keeps tags in memory, so you should not use
// any unbounded value in the tag value.
func New() MetricsClient {
	client := MetricsClient{
		counters: map[string]*prometheus.CounterVec{
			"ingress_topic_counter": promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "ingress_topic_counter",
				Help: "A counter for incoming messages in a kafka topic",
			}, []string{"topic"}),
			"egress_topic_counter": promauto.NewCounterVec(prometheus.CounterOpts{
				Name: "egress_topic_counter",
				Help: "A counter for outgoing canonical transactions",
			}, []string{"topic", "status"}),
		},
		histograms: map[string]*prometheus.HistogramVec{
			"transaction_validation_time_ms": promauto.NewHistogramVec(prometheus.HistogramOpts{
				Name: "transaction_validation_time_ms",
				Help: "Time to consume, validate, and write back transaction",
				// this list represents the tags name for this metric
			}, []string{"topic", "status"}),
		},
	}

	return client
}

func (m MetricsClient) Histogram(name string, value float64, tags []string) {
	if h, ok := m.histograms[name]; ok {
		h.WithLabelValues(tags...).Observe(value)

		return
	}
	// TODO: log if metric not found
}

func (m MetricsClient) Count(name string, value int64, tags []string) {
	if h, ok := m.counters[name]; ok {
		h.WithLabelValues(tags...).Add(float64(value))

		return
	}
	// TODO: log if metric not found
}
