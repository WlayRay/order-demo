package main

import (
	"context"
	"github.com/WlayRay/order-demo/common/broker"
	_ "github.com/WlayRay/order-demo/common/config"
	"github.com/WlayRay/order-demo/common/logging"
	"github.com/WlayRay/order-demo/common/server"
	"github.com/WlayRay/order-demo/common/tracing"

	"github.com/WlayRay/order-demo/payment/infrastructure/consumer"
	"github.com/WlayRay/order-demo/payment/service"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	logging.Init()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceName := viper.GetString("payment.service-name")
	serverType := viper.GetString("payment.server-to-run")

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName, viper.GetFloat64("jaeger.sampling-rate"))
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	// 初始化消息队列
	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	paymentHandler := NewPaymentHandler(ch)
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		zap.L().Panic("unsupported server type", zap.String("serverType", serverType))
	default:
		zap.L().Panic("unknown server type", zap.String("serverType", serverType))
	}
}
