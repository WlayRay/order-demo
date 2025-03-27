package main

import (
	"context"
	"github.com/WlayRay/order-demo/common/config"
	"github.com/WlayRay/order-demo/common/discovery"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/common/logging"
	"github.com/WlayRay/order-demo/common/server"
	"github.com/WlayRay/order-demo/order/ports"
	"github.com/WlayRay/order-demo/order/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		zap.L().Fatal("init config error", zap.Error(err))
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		orderpb.RegisterOrderServiceServer(server, ports.NewGRPCServer(application))
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})
}
