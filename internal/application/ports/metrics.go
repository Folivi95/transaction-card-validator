//go:generate moq -out mocks/metrics_client_moq.go -pkg=mocks . MetricsClient

package ports

// MetricsClient sends various kinds of metrics to monitoring tools.
type MetricsClient interface {
	Histogram(name string, value float64, tags []string)
	Count(name string, value int64, tags []string)
}

type DummyMetricsClient struct{}

func (DummyMetricsClient) Histogram(_ string, _ float64, _ []string) {}
func (DummyMetricsClient) Count(_ string, _ int64, _ []string)       {}
