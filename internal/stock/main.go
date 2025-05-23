package main

import (
	"context"
	_ "github.com/WlayRay/order-demo/common/config"
	"github.com/WlayRay/order-demo/common/discovery"
	"github.com/WlayRay/order-demo/common/genproto/stockpb"
	"github.com/WlayRay/order-demo/common/logging"
	"github.com/WlayRay/order-demo/common/server"
	"github.com/WlayRay/order-demo/common/tracing"
	"github.com/WlayRay/order-demo/stock/ports"
	"github.com/WlayRay/order-demo/stock/service"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName, viper.GetFloat64("jaeger.sampling-rate"))
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	application := service.NewApplication(ctx)

	UnRegisterFunc, err := discovery.RegisterToETCD(ctx, serviceName)
	if err != nil {
		zap.L().Fatal("register to etcd error", zap.Error(err))
	}
	defer func() {
		if err := UnRegisterFunc(); err != nil {
			_ = UnRegisterFunc()
		}
	}()

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			stockpb.RegisterStockServiceServer(server, ports.NewGRPCServer(application))
		})
	case "http":
		// DoNothing
	default:
		panic("unknown server type")
	}
}
