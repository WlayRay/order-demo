package grpcClient

import (
	"context"
	"github.com/WlayRay/order-demo/common/genproto/stockpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (stockpb.StockServiceClient, func(), error) {
	grpcAddr := viper.GetString("stock.grpc-addr")
	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, nil, err
	}
	conn, newClientErr := grpc.NewClient(grpcAddr, opts...)
	if newClientErr != nil {
		return nil, nil, newClientErr
	}
	return stockpb.NewStockServiceClient(conn), func() {
		_ = conn.Close()
	}, nil
}

func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, nil
}
