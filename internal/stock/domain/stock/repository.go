package stock

import (
	"context"
	"fmt"
	"github.com/WlayRay/order-demo/stock/entity"
	"strings"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*entity.Item, error)
}

type NotFoundError struct {
	Missing []string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("not found items: %v", strings.Join(n.Missing, ","))
}
