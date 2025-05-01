package command

import (
	"context"

	"github.com/WlayRay/order-demo/common/decorator"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/payment/domain"
	"go.uber.org/zap"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	processor domain.Processor
	orderGRPC OrderService
}

func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}

	zap.L().Info("CreatePaymentHandler.Handle: CreatePaymentLink", zap.Any("order", cmd.Order), zap.Any("link", link))
	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting_for_payment", // 更新订单状态为待支付
		PaymentLink: link,                  // 创建支付链接
		Items:       cmd.Order.Items,
	}
	err = c.orderGRPC.UpdateOrder(ctx, newOrder) // 重新调用order服务的UpdateOrder方法，更新订单信息
	if err != nil {
		return "", err
	}
	return link, nil
}

func NewCreatePaymentHandler(processor domain.Processor, orderGRPC OrderService, logger *zap.Logger, metricClient decorator.MetricsClient) CreatePaymentHandler {
	if processor == nil {
		panic("processor is nil")
	}

	return decorator.ApplyCommandDecorators(
		createPaymentHandler{
			processor: processor,
			orderGRPC: orderGRPC,
		},
		logger,
		metricClient,
	)
}
