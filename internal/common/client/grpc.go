package grpcClient

import (
	"context"
	"errors"
	"github.com/WlayRay/order-demo/common/discovery"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/common/genproto/stockpb"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"time"
)

func NewStockGRPCClient(ctx context.Context) (stockpb.StockServiceClient, func() error, error) {
	if !WaitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout")) {
		return nil, nil, errors.New("stock grpc not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	if grpcAddr == "" {
		zap.L().Warn("stock service not found")
	}

	opts := grpcDialOpts(grpcAddr)
	conn, newClientErr := grpc.NewClient(grpcAddr, opts...)
	if newClientErr != nil {
		return nil, nil, newClientErr
	}
	return stockpb.NewStockServiceClient(conn), func() error {
		return conn.Close()
	}, nil
}

func NewOrderGRPCClient(ctx context.Context) (orderpb.OrderServiceClient, func() error, error) {
	if !WaitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout")) {
		return nil, nil, errors.New("order grpc not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	if grpcAddr == "" {
		zap.L().Warn("stock service not found")
	}

	opts := grpcDialOpts(grpcAddr)
	conn, newClientErr := grpc.NewClient(grpcAddr, opts...)
	if newClientErr != nil {
		return nil, nil, newClientErr
	}
	return orderpb.NewOrderServiceClient(conn), func() error {
		return conn.Close()
	}, nil
}

func grpcDialOpts(_ string) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
}

func WaitForOrderGRPCClient(timeout time.Duration) bool {
	zap.L().Info("waiting for order service", zap.Duration("timeout", timeout))
	return waitFor(viper.GetString("order.grpc-addr"), timeout)
}

func WaitForStockGRPCClient(timeout time.Duration) bool {
	zap.L().Info("waiting for stock service", zap.Duration("timeout", timeout))
	return waitFor(viper.GetString("stock.grpc-addr"), timeout)
}

func waitFor(addr string, timeout time.Duration) bool {
	portAvailable := make(chan struct{}, 1)
	timeoutCh := time.After(timeout)

	go func() {
		defer close(portAvailable)
		for {
			select {
			case <-timeoutCh:
				zap.L().Fatal("timeout waiting for grpc server", zap.String("addr", addr))
				return
			default:
				conn, err := net.Dial("tcp", addr)
				if err == nil {
					_ = conn.Close()
					portAvailable <- struct{}{}
					return
				}
				time.Sleep(200 * time.Millisecond)
			}

		}
	}()

	for {
		select {
		case <-portAvailable:
			return true
		case <-timeoutCh:
			return false
		}
	}
}
