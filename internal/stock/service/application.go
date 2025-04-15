package service

import (
	"context"
	_ "github.com/WlayRay/order-demo/common/config"
	"github.com/WlayRay/order-demo/common/metrics"
	"github.com/WlayRay/order-demo/stock/adapters"
	"github.com/WlayRay/order-demo/stock/app"
	"github.com/WlayRay/order-demo/stock/app/query"
	"github.com/WlayRay/order-demo/stock/infrastructure/integration"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewApplication(ctx context.Context) app.Application {
	stockRepo := adapters.NewStockRepositoryPG(adapters.NewEntClient())
	stripeAPI := integration.NewStripeAPI()
	logger := zap.L()
	metricClient := metrics.NewPrometheusMetricsClient(
		&metrics.PrometheusMetricsClientConfig{
			Host:        viper.GetString("stock.metrics-export-addr"),
			ServiceName: viper.GetString("stock.service-name"),
		})
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logger, metricClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricClient),
		},
	}
}
