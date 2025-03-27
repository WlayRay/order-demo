package domain

import (
	"context"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error)
}
