package metrics

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	path            string
	port            string
	server          *http.Server
	serverMu        sync.Mutex
	isServerRunning bool
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	path = os.Getenv("PROMETHEUS_PATH")
	port = os.Getenv("PROMETHEUS_ADDR")
	if path == "" || port == "" {
		panic(errors.New("Prometheus读取环境变量失败"))
	}
}

func StartPrometheusServer() error {
	serverMu.Lock()
	defer serverMu.Unlock()

	// 防止重复启动
	if isServerRunning {
		errors.New("Prometheus 服务已经在运行")
	}

	mux := http.NewServeMux()
	mux.Handle(path, promhttp.Handler())

	server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(errors.Wrap(err, "Prometheus 监控服务启动失败"))
		}
	}()

	isServerRunning = true
	return nil
}

// StopPrometheusServer 优雅关闭 Prometheus 服务
func StopPrometheusServer() error {
	serverMu.Lock()
	defer serverMu.Unlock()

	if !isServerRunning || server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "Prometheus 服务关闭失败")
	}

	isServerRunning = false
	server = nil
	return nil
}
