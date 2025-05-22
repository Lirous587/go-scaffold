package metrics

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
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
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"action", "status"},
	)
	h := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.5, 1, 2, 5},
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
	port string
)

func init() {
	_ = godotenv.Load()
	path = os.Getenv("PROMETHEUS_PATH")
	port = os.Getenv("PROMETHEUS_ADDR")
	if path == "" || port == "" {
		panic(errors.New("Prometheus读取环境变量失败"))
	}
}

func StartPrometheusServer() {
	http.Handle(path, promhttp.Handler())
	go http.ListenAndServe(":"+port, nil)
}
