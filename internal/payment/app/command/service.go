package command

import (
	"context"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
