package service

import (
	"context"
	grpcClient "github.com/WlayRay/order-demo/common/client"
	_ "github.com/WlayRay/order-demo/common/config"
	"github.com/WlayRay/order-demo/common/metrics"
	"github.com/WlayRay/order-demo/payment/adapters"
	"github.com/WlayRay/order-demo/payment/app"
	"github.com/WlayRay/order-demo/payment/app/command"
	"github.com/WlayRay/order-demo/payment/infrastructure/processor"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err.Error())
	}

	orderGRPC := adapters.NewOrderGRPC(orderClient)
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))

	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor *processor.StripeProcessor) app.Application {
	logger := zap.L()
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}

}
