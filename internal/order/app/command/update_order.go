package command

import (
	"context"
	"github.com/WlayRay/order-demo/common/decorator"
	domain "github.com/WlayRay/order-demo/order/domain/order"
	"go.uber.org/zap"
)

type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, any]

type updateOrderHandler struct {
	orderRepo domain.Repository
	//stockGRPC
}

func NewUpdateOrderHandler(orderRepo domain.Repository, logger *zap.Logger, metricClient decorator.MetricsClient) UpdateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}

	return decorator.ApplyCommandDecorators[UpdateOrder, any](
		updateOrderHandler{
			orderRepo: orderRepo,
		},
		logger,
		metricClient,
	)
}

func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (any, error) {
	if cmd.UpdateFn != nil {
		zap.L().Warn("UpdateOrderHandler.Handle: UpdateFn is nil", zap.Any("order:", cmd.Order))
		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		}
	}
	if err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
		return nil, err
	}
	return nil, nil
}
