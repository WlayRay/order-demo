package service

import (
	"context"
	"github.com/WlayRay/order-demo/common/broker"
	grpcClient "github.com/WlayRay/order-demo/common/client"
	"github.com/WlayRay/order-demo/common/config"
	"github.com/WlayRay/order-demo/common/metrics"
	"github.com/WlayRay/order-demo/order/adapters"
	"github.com/WlayRay/order-demo/order/adapters/grpc"
	"github.com/WlayRay/order-demo/order/app" // 注意这里是order
	"github.com/WlayRay/order-demo/order/app/command"
	"github.com/WlayRay/order-demo/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		zap.L().Fatal("init config error", zap.Error(err))
	}
}

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeFn, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)

	// 初始化消息队列
	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)

	return newApplication(ctx, stockGRPC, ch), func() {
		_ = closeFn()
		_ = ch.Close()
		_ = closeCh()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := zap.L()
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
