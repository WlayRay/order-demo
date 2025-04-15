package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

type PrometheusMetricsClient struct {
	registry *prometheus.Registry
}

var dynamicCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "dynamic_counter",
		Help: "count_custom_keys",
	}, []string{"key"})

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
	p.registry.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// custom collector
	if err := p.registry.Register(dynamicCounter); err != nil {
		panic(err)
	}

	// metadata wrap
	prometheus.WrapRegistererWith(prometheus.Labels{"serviceName": conf.ServiceName}, p.registry)

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
	dynamicCounter.WithLabelValues(key).Add(float64(value))
}
