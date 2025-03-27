package discovery

import (
	"context"
	"github.com/WlayRay/order-demo/common/discovery/etcd"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

func RegisterToETCD(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := etcd.GetEtcdClient(viper.GetStringSlice("etcd.endpoints"))
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
	zap.L().Info("registered to consul",
		zap.String("serviceName", serviceName),
		zap.String("instanceID", instanceID))

	return func() error {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return registry.Unregister(timeoutCtx, instanceID, serviceName)
	}, nil
}
