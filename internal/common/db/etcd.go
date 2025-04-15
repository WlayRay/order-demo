package db

import (
	"github.com/spf13/viper"
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
func GetEtcdClient() (*etcdv3.Client, error) {
	once.Do(func() {
		client, err = etcdv3.New(etcdv3.Config{
			Endpoints:          viper.GetStringSlice("etcd.endpoints"),
			DialTimeout:        5 * time.Second,
			MaxCallSendMsgSize: 10 * 1024 * 1024, // 设置发送消息的最大大小
			MaxCallRecvMsgSize: 10 * 1024 * 1024, // 设置接收消息的最大大小
		})
	})

	return client, err
}
