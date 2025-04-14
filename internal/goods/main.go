package main

import (
	"context"
	"github.com/WlayRay/order-demo/common/broker"
	grpcClient "github.com/WlayRay/order-demo/common/client"
	_ "github.com/WlayRay/order-demo/common/config"
	"github.com/WlayRay/order-demo/common/logging"
	"github.com/WlayRay/order-demo/goods/adaptors"
	"github.com/WlayRay/order-demo/goods/infrastructure/consumer"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"github.com/WlayRay/order-demo/common/tracing"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("goods.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	orderClient, closeFunc, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		zap.L().Fatal("Failed to create order client", zap.Error(err))
	}
	defer func() { _ = closeFunc() }()

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

	orderGRPC := adaptors.NewOrderGRPC(orderClient)
	go consumer.NewConsumer(orderGRPC).Listen(ch)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalCh:
			zap.L().Info("Received shutdown signal, shutting down gracefully...")
			os.Exit(0)
		}
	}
}
