package grpcClient

import (
	"context"
	"github.com/WlayRay/order-demo/common/discovery"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/common/genproto/stockpb"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (stockpb.StockServiceClient, func() error, error) {
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	if grpcAddr == "" {
		zap.L().Warn("stock service not found")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, nil, err
	}
	conn, newClientErr := grpc.NewClient(grpcAddr, opts...)
	if newClientErr != nil {
		return nil, nil, newClientErr
	}
	return stockpb.NewStockServiceClient(conn), func() error {
		return conn.Close()
	}, nil
}

func NewOrderGRPCClient(ctx context.Context) (orderpb.OrderServiceClient, func() error, error) {
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	if grpcAddr == "" {
		zap.L().Warn("stock service not found")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, nil, err
	}
	conn, newClientErr := grpc.NewClient(grpcAddr, opts...)
	if newClientErr != nil {
		return nil, nil, newClientErr
	}
	return orderpb.NewOrderServiceClient(conn), func() error {
		return conn.Close()
	}, nil
}

func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, nil
}
