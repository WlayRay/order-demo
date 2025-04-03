package etcd

import (
	"context"
	"errors"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

type Registry struct {
	client *etcdv3.Client
}

// GetEtcdClient returns a new etcd registry client.
func GetEtcdClient(etcdEndpoints []string) (*Registry, error) {
	once.Do(func() {
		cli, e := etcdv3.New(etcdv3.Config{
			Endpoints:          etcdEndpoints,
			DialTimeout:        5 * time.Second,
			MaxCallSendMsgSize: 10 * 1024 * 1024, // 设置发送消息的最大大小
			MaxCallRecvMsgSize: 10 * 1024 * 1024, // 设置接收消息的最大大小
		})

		if e != nil {
			err = e
			return
		}

		etcdClient = &Registry{
			client: cli,
		}
	})

	if err != nil {
		return nil, err
	}
	return etcdClient, nil
}

// Register registers a service instance in etcd.
func (r Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid hostPort")
	}

	key := "/" + serviceName + "/" + instanceID + "/" + hostPort

	// 创建3秒的租约
	lease, err := r.client.Grant(ctx, 3)
	if err != nil {
		return err
	}
	// 服务注册(向ETCD中写入一个key)
	_, err = r.client.Put(ctx, key, "", etcdv3.WithLease(lease.ID))
	return err
}

// Unregister unregisters a service instance from etcd.
func (r Registry) Unregister(ctx context.Context, instanceID, serviceName string) error {
	zap.L().Info("unregister service",
		zap.String("serviceName", serviceName),
		zap.String("instanceID", instanceID))
	_, err := r.client.Delete(ctx, "/"+serviceName+"/"+instanceID, etcdv3.WithPrefix())
	return err
}

// Discover discovers service instances from etcd.
func (r Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	resp, err := r.client.Get(ctx, "/"+serviceName+"/", etcdv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	addresses := make([]string, 0, resp.Count)
	for _, kv := range resp.Kvs {
		ips := strings.Split(string(kv.Key), "/")
		addresses = append(addresses, ips[(len(ips)-1)])
	}
	return addresses, nil
}

// HealthCheck checks the health of a service instance in etcd.
func (r Registry) HealthCheck(instanceID, serviceName string) error {
	ctx := context.Background()

	// 获取服务实例的所有键
	prefix := "/" + serviceName + "/" + instanceID
	resp, err := r.client.Get(ctx, prefix, etcdv3.WithPrefix())
	if err != nil {
		return err
	}

	if resp.Count == 0 {
		return errors.New("service instance not found")
	}

	// 为每个键续租
	for _, kv := range resp.Kvs {
		leaseResp, err := r.client.Grant(ctx, 3) // 创建3秒的租约
		if err != nil {
			return err
		}

		// 更新键值对的租约
		_, err = r.client.Put(ctx, string(kv.Key), string(kv.Value), etcdv3.WithLease(leaseResp.ID))
		if err != nil {
			return err
		}

		// 自动续租
		keepAliveChan, keepAliveErr := r.client.KeepAlive(ctx, leaseResp.ID)
		if keepAliveErr != nil {
			return keepAliveErr
		}

		// 处理续租响应
		go func() {
			for resp := range keepAliveChan {
				if resp == nil {
					zap.L().Warn("keep alive channel closed",
						zap.String("serviceName", serviceName),
						zap.String("instanceID", instanceID))
					return
				}
				/*zap.L().Debug("keep alive success",
				zap.String("serviceName", serviceName),
				zap.String("instanceID", instanceID),
				zap.Int64("ttl", resp.TTL))*/
			}
		}()
	}

	return nil
}

var (
	etcdClient *Registry
	once       sync.Once
	err        error
)
