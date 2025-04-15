package discovery

import (
	"context"
	"fmt"
	"github.com/WlayRay/order-demo/common/discovery/etcd"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

// RegisterToETCD registers a service to etcd.
func RegisterToETCD(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := etcd.GetRegistry()
	if err != nil {
		return func() error {
			return nil
		}, err
	}

	instanceID := GenerateInstanceID(serviceName)
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")

	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		return func() error {
			return nil
		}, err
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				zap.L().Panic("health check failed", zap.Error(err))
			}
			time.Sleep(time.Second * 2)
		}
	}()
	zap.L().Info("registered to etcd",
		zap.String("serviceName", serviceName),
		zap.String("instanceID", instanceID))

	return func() error {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return registry.Unregister(timeoutCtx, instanceID, serviceName)
	}, nil
}

// GetServiceAddr retrieves the address of a service from etcd.
func GetServiceAddr(ctx context.Context, serviceName string) (string, error) {
	registry, err := etcd.GetRegistry()
	if err != nil {
		return "", err
	}
	addresses, discoverErr := registry.Discover(ctx, serviceName)
	if discoverErr != nil {
		return "", discoverErr
	}
	if len(addresses) == 0 {
		return "", fmt.Errorf("%s no alive service found", serviceName)
	}
	i := rand.Intn(len(addresses))
	zap.L().Info("get service addr",
		zap.String("serviceName", serviceName),
		zap.String("addr", addresses[i]))
	return addresses[i], nil
}
