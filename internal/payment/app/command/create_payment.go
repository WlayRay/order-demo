package command

import (
	"context"
	"github.com/WlayRay/order-demo/common/decorator"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/payment/domain"
	"go.uber.org/zap"
)

type CreatePaymentLink struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePaymentLink, string]

type createPaymentHandler struct {
	processor domain.Processor
	orderGRPC OrderService
}

func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePaymentLink) (string, error) {
	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}

	zap.L().Info("CreatePaymentHandler.Handle: CreatePaymentLink", zap.Any("order:", cmd.Order.ID), zap.Any("link:", link))
	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting_for_payment",
		PaymentLink: link,
		Items:       cmd.Order.Items,
	}
	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
	if err != nil {
		return "", err
	}
	return link, nil
}

func NewCreatePaymentHandler(processor domain.Processor, orderGRPC OrderService, logger *zap.Logger, metricClient decorator.MetricsClient) CreatePaymentHandler {
	if processor == nil {
		panic("processor is nil")
	}

	return decorator.ApplyCommandDecorators[CreatePaymentLink, string](
		createPaymentHandler{
			processor: processor,
			orderGRPC: orderGRPC,
		},
		logger,
		metricClient,
	)
}
