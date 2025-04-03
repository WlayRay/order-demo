package discovery

import (
	"context"
	"fmt"
	"math/rand"
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
	x := rand.New(rand.NewSource(time.Now().Unix()))
	return fmt.Sprintf("%s-%d", serviceName, x.Uint64())
}
