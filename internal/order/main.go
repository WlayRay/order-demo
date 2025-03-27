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
	"time"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		zap.L().Fatal("init config error", zap.Error(err))
	}
}

func main() {
	//fmt.Println("环境变量:", os.Getenv("STRIPE_KEY")) // 验证环境变量是否生效
	//zap.L().Fatal(viper.GetString("stripe-key"))
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
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
