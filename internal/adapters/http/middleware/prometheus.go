package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap"
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_duration_seconds",
	Help: "Duration of HTTP requests.",
	Buckets: buckets(
		1*time.Millisecond,
		10*time.Millisecond,
		100*time.Millisecond,
		250*time.Millisecond,
		500*time.Millisecond,
		1*time.Second,
		2500*time.Millisecond,
		5*time.Second,
		10*time.Second,
		30*time.Second,
	),
}, []string{"handler", "method", "code"})

// Prometheus provides a middleware that will measure and publish metrics for each http request.
// Provides information about the path(template), http method and http response code.
func Prometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		// using the template instead of the path to keep cardinality low, meaning not use the real parameters to avoid having multiple
		// paths
		path, err := route.GetPathTemplate()
		if err != nil {
			// handle the error but don't fail the request just because we can't measure it
			zapctx.Error(r.Context(), "error extracting path", zap.Error(err))
			next.ServeHTTP(w, r)

			return
		}

		promhttp.
			InstrumentHandlerDuration(httpDuration.MustCurryWith(
				prometheus.Labels{"handler": path}),
				next).
			ServeHTTP(w, r)
	})
}

func buckets(durations ...time.Duration) []float64 {
	s := make([]float64, len(durations))
	for i, d := range durations {
		s[i] = d.Seconds()
	}

	return s
}
