package db

import (
	etcdv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var (
	client *etcdv3.Client
	once   sync.Once
	err    error
)

// GetEtcdClient returns a new etcd registry client.
func GetEtcdClient(etcdEndpoints []string) (*etcdv3.Client, error) {
	once.Do(func() {
		client, err = etcdv3.New(etcdv3.Config{
			Endpoints:          etcdEndpoints,
			DialTimeout:        5 * time.Second,
			MaxCallSendMsgSize: 10 * 1024 * 1024, // 设置发送消息的最大大小
			MaxCallRecvMsgSize: 10 * 1024 * 1024, // 设置接收消息的最大大小
		})
	})

	return client, err
}
