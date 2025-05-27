package main

import (
	"context"
	"github.com/WlayRay/order-demo/common/broker"
	_ "github.com/WlayRay/order-demo/common/config"
	"github.com/WlayRay/order-demo/common/discovery"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/common/logging"
	"github.com/WlayRay/order-demo/common/server"
	"github.com/WlayRay/order-demo/common/tracing"
	"github.com/WlayRay/order-demo/order/infrastructure/consumer"
	"github.com/WlayRay/order-demo/order/ports"
	"github.com/WlayRay/order-demo/order/service"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName, viper.GetFloat64("jaeger.sampling-rate"))
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	UnRegisterFunc, err := discovery.RegisterToETCD(ctx, serviceName)
	if err != nil {
		zap.L().Fatal("register to etcd error", zap.Error(err))
	}
	defer func() {
		if err := UnRegisterFunc(); err != nil {
			_ = UnRegisterFunc()
		}
	}()

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

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		orderpb.RegisterOrderServiceServer(server, ports.NewGRPCServer(application))
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		if viper.GetBool("enable-profiling") {
			pprof.Register(router)
		}

		router.StaticFile("/success", "..\\..\\public\\success.html")
		ports.RegisterHandlersWithOptions(router, HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})
}
