package service

import (
	"context"
	grpcClient "github.com/WlayRay/order-demo/common/client"
	"github.com/WlayRay/order-demo/common/metrics"
	"github.com/WlayRay/order-demo/order/adapters"
	"github.com/WlayRay/order-demo/order/adapters/grpc"
	"github.com/WlayRay/order-demo/order/app" // 注意这里是order
	"github.com/WlayRay/order-demo/order/app/command"
	"github.com/WlayRay/order-demo/order/app/query"
	"go.uber.org/zap"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeFn, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC), func() {
		closeFn()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := zap.L()
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
