package discovery

import (
	"context"
	"fmt"
	"github.com/WlayRay/order-demo/common/lib"
	"hash/fnv"
	"time"
)

type Registry interface {
	Register(ctx context.Context, instanceID, serviceName, hostPort string) error
	Unregister(ctx context.Context, instanceID, serviceName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceID, serviceName string) error
}

// GenerateInstanceID generates a unique instance ID for a service.
func GenerateInstanceID(serviceName string) string {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()
	ip, err := lib.GetLocalIP()

	h := fnv.New64a()
	_, _ = h.Write([]byte(ip))
	snowflakeInstance, err := lib.GetSnowflakeInstance(h.Sum64()%1024, 10*time.Millisecond)

	id, err := snowflakeInstance.GetID()
	return fmt.Sprintf("%s-%d", serviceName, id)
}
