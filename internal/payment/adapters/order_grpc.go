package adapters

import (
	"context"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/common/tracing"
	"go.uber.org/zap"
)

type OrderGRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) error {
	ctx, span := tracing.Start(ctx, "OrderGRPC.UpdateOrder")
	defer span.End()

	_, err := o.client.UpdateOrder(ctx, order)
	if err != nil {
		zap.L().Warn("failed to update order", zap.Error(err))
	}
	return err
}
