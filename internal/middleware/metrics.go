package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics contains application HTTP metrics and an isolated Prometheus registry.
type Metrics struct {
	requests *prometheus.CounterVec
	latency  *prometheus.HistogramVec
	inFlight prometheus.Gauge
	registry *prometheus.Registry
}

func NewMetrics() *Metrics {
	m := &Metrics{
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total", Help: "Total HTTP requests.",
		}, []string{"method", "path", "status"}),
		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_request_duration_seconds", Help: "HTTP request duration in seconds.",
		}, []string{"method", "path"}),
		inFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "http_requests_in_flight", Help: "HTTP requests currently being served.",
		}),
		registry: prometheus.NewRegistry(),
	}
	m.registry.MustRegister(m.requests, m.latency, m.inFlight)
	return m
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.inFlight.Inc()
		started := time.Now()
		c.Next()
		m.inFlight.Dec()
		path := c.FullPath()
		if path == "" {
			path = "unmatched"
		}
		method := c.Request.Method
		m.requests.WithLabelValues(method, path, strconv.Itoa(c.Writer.Status())).Inc()
		m.latency.WithLabelValues(method, path).Observe(time.Since(started).Seconds())
	}
}

func (m *Metrics) Handler() gin.HandlerFunc {
	handler := promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
	return gin.WrapH(handler)
}
