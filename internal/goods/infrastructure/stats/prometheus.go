package stats

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type GoodsMetrics struct {
	salesVolumeCounter *prometheus.CounterVec
}

type PrometheusStats struct {
	registry    *prometheus.Registry
	host        string
	serviceName string
	*GoodsMetrics
}

func NewPrometheusStats(host string, serviceName string) *PrometheusStats {
	return &PrometheusStats{
		registry:     prometheus.NewRegistry(),
		host:         host,
		serviceName:  serviceName,
		GoodsMetrics: newGoodsMetrics(),
	}
}

func (p *PrometheusStats) Start() error {
	wr := prometheus.WrapRegistererWith(
		prometheus.Labels{"serviceName": p.serviceName},
		p.registry,
	)

	wr.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		p.GoodsMetrics.salesVolumeCounter,
	)

	// export
	http.Handle("/metrics", promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}))
	go func() {
		zap.L().Info("Starting metrics server", zap.String("host", p.host))
		if err := http.ListenAndServe(p.host, nil); err != nil {
			zap.L().Fatal("Failed to start metrics server", zap.Error(err))
		}
	}()

	return nil
}

func newGoodsMetrics() *GoodsMetrics {
	goodsMetrics := &GoodsMetrics{}
	goodsMetrics.salesVolumeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goods_counter",
			Help: "产品销售量",
		},
		[]string{"name", "category", "brand"},
	)
	return goodsMetrics
}

func (g *GoodsMetrics) IncSalesVolume(goodName, category, brand string, value int) {
	g.salesVolumeCounter.WithLabelValues(goodName, category, brand).Add(float64(value))
}
