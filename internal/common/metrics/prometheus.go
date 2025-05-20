package metrics

import (
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

type PrometheusClient struct {
	counter   *prometheus.CounterVec
	histogram *prometheus.HistogramVec
}

func NewPrometheusClient() *PrometheusClient {
	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_command_metrics",
			Help: "Command metrics",
		},
		[]string{"action", "status"},
	)
	h := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_command_duration_seconds",
			Help:    "Command duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"action", "status"},
	)
	prometheus.MustRegister(c, h)
	return &PrometheusClient{counter: c, histogram: h}
}

func (p *PrometheusClient) Inc(action, status string, value int) {
	p.counter.WithLabelValues(action, status).Add(float64(value))
}

func (p *PrometheusClient) ObserveDuration(action, status string, seconds float64) {
	p.histogram.WithLabelValues(action, status).Observe(seconds)
}

var (
	path string
	addr string
)

func init() {
	_ = godotenv.Load()
	path = os.Getenv("")
	path = os.Getenv("")

}

func StartPrometheusServer(path, addr string) {
	http.Handle(path, promhttp.Handler())
	go http.ListenAndServe(addr, nil)
}
