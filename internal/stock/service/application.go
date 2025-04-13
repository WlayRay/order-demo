package service

import (
	"context"
	"github.com/WlayRay/order-demo/common/metrics"
	"github.com/WlayRay/order-demo/stock/adapters"
	"github.com/WlayRay/order-demo/stock/app"
	"github.com/WlayRay/order-demo/stock/app/query"
	"github.com/WlayRay/order-demo/stock/infrastructure/integration"
	"go.uber.org/zap"
)

func NewApplication(ctx context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	stripeAPI := integration.NewStripeAPI()
	logger := zap.L()
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logger, metricClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricClient),
		},
	}
}
