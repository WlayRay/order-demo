package stock

import (
	"context"
	"fmt"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"strings"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
}

type NotFoundError struct {
	Missing []string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("not found items: %v", strings.Join(n.Missing, ","))
}
