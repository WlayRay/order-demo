package command

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/WlayRay/order-demo/common/broker"
	"github.com/WlayRay/order-demo/common/decorator"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/order/app/query"
	domain "github.com/WlayRay/order-demo/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type CreateOrder struct {
	CustomerID string
	Items      []*orderpb.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	chanel    *amqp.Channel
}

func NewCreateOrderHandler(orderRepo domain.Repository, stockGRPC query.StockService, chanel *amqp.Channel, logger *zap.Logger, metricClient decorator.MetricsClient) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	if stockGRPC == nil {
		panic("stockGRPC is nil")
	}
	if chanel == nil {
		panic("chanel is nil")
	}

	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{
			orderRepo: orderRepo,
			stockGRPC: stockGRPC,
			chanel:    chanel,
		},
		logger,
		metricClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	validItems, err := c.validate(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}

	o, createErr := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID:  cmd.CustomerID,
		Items:       validItems,
		PaymentLink: "price_1R7HVgPNegMNE0WfuwRkVr6b",
	})
	if createErr != nil {
		return nil, createErr
	}

	q, queueErr := c.chanel.QueueDeclare(broker.EventOrderCreated, false, true, false, false, nil)
	if queueErr != nil {
		return nil, queueErr
	}

	marshalOrder, jsonErr := json.Marshal(o)
	if jsonErr != nil {
		return nil, jsonErr
	}

	err = c.chanel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         marshalOrder,
	})
	if err != nil {
		return nil, err
	}

	return &CreateOrderResult{OrderID: o.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}
	items = packItems(items)
	resp, err := c.stockGRPC.CheckItemsInStock(items, ctx)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}
	var res []*orderpb.ItemWithQuantity
	for id, quantity := range merged {
		res = append(res, &orderpb.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}
	return res
}
