package service

import (
	"context"
	"github.com/WlayRay/order-demo/common/metrics"
	"github.com/WlayRay/order-demo/order/adapters"
	"github.com/WlayRay/order-demo/order/app" // 注意这里是order
	"github.com/WlayRay/order-demo/order/app/query"
	"go.uber.org/zap"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := zap.L()
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
