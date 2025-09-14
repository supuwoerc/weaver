package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// http请求总数统计(计数器) promauto自动注册
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"method", "path", "status"})

	// http请求耗时统计(直方图) promauto自动注册
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 1, 1.5, 2, 3},
	}, []string{"method", "path"})
	// 活跃连接(仪表) promauto自动注册
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)
)

type PrometheusMiddleware struct {
}

func NewPrometheusMiddleware() *PrometheusMiddleware {
	return &PrometheusMiddleware{}
}

func (r *PrometheusMiddleware) Prometheus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// active connections
		activeConnections.Inc()
		defer activeConnections.Dec()
		path := ctx.FullPath()
		// request duration observer
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(ctx.Request.Method, path))
		ctx.Next()
		status := ctx.Writer.Status()
		// http request count +1
		httpRequestsTotal.WithLabelValues(ctx.Request.Method, path, http.StatusText(status)).Inc()
		// http request duration
		timer.ObserveDuration()
	}
}
