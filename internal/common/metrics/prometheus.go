package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type PrometheusMetricsClient struct {
	registry *prometheus.Registry
}

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dynamic_counter",
			Help: "The number of successful and failed execution of business methods",
		}, []string{"key"})

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "dynamic_duration",
			Help: "The duration of successful and failed execution of business methods",
		}, []string{"key"})
)

type PrometheusMetricsClientConfig struct {
	Host        string
	ServiceName string
}

func NewPrometheusMetricsClient(config *PrometheusMetricsClientConfig) *PrometheusMetricsClient {
	client := &PrometheusMetricsClient{}
	client.initPrometheus(config)
	return client
}

func (p PrometheusMetricsClient) initPrometheus(conf *PrometheusMetricsClientConfig) {
	p.registry = prometheus.NewRegistry()

	wr := prometheus.WrapRegistererWith(
		prometheus.Labels{"serviceName": conf.ServiceName},
		p.registry,
	)

	wr.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	if err := wr.Register(requestCounter); err != nil {
		panic(err)
	}
	if err := wr.Register(requestDuration); err != nil {
		panic(err)
	}

	// export
	http.Handle("/metrics", promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}))
	go func() {
		zap.L().Info("Starting metrics server", zap.String("host", conf.Host))
		if err := http.ListenAndServe(conf.Host, nil); err != nil {
			zap.L().Fatal("Failed to start metrics server", zap.Error(err))
		}
	}()
}

func (p PrometheusMetricsClient) Inc(key string, value int) {
	requestCounter.WithLabelValues(key).Add(float64(value))
}

func (p PrometheusMetricsClient) Observe(key string, value float64) {
	requestDuration.WithLabelValues(key).Observe(value)
}
