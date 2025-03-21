package main

import (
	"github.com/WlayRay/order-demo/common/genproto/stockpb"
	"github.com/WlayRay/order-demo/common/server"
	"github.com/WlayRay/order-demo/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("order.service-name")
	serverType := viper.GetString("order.server-to-run")
	switch serverType {
	case "gprc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			stockpb.RegisterStockServiceServer(server, ports.NewGRPCServer())
		})
	case "http":
		// DoNothing
	default:
		panic("unknown server type")
	}
}
