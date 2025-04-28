package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WlayRay/order-demo/common/broker"
	"github.com/WlayRay/order-demo/common/decorator"
	"github.com/WlayRay/order-demo/order/app/query"
	"github.com/WlayRay/order-demo/order/convertor"
	domain "github.com/WlayRay/order-demo/order/domain/order"
	"github.com/WlayRay/order-demo/order/entity"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type CreateOrder struct {
	CustomerID string
	Items      []*entity.ItemWithQuantity
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

	return decorator.ApplyCommandDecorators(
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
	q, queueErr := c.chanel.QueueDeclare(broker.EventOrderCreated, false, true, false, false, nil)
	if queueErr != nil {
		return nil, queueErr
	}

	t := otel.Tracer("rabbitmq")
	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", q.Name))
	defer span.End()

	validItems, err := c.validate(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}

	o, createErr := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      validItems,
	})
	if createErr != nil {
		return nil, createErr
	}

	marshalOrder, jsonErr := json.Marshal(o)
	if jsonErr != nil {
		return nil, jsonErr
	}

	header := broker.InjectRabbitMQHeaders(ctx)
	err = c.chanel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         marshalOrder,
		Headers:      header,
	})
	if err != nil {
		return nil, err
	}

	return &CreateOrderResult{OrderID: o.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}
	items = packItems(items)
	resp, err := c.stockGRPC.CheckItemsInStock(ctx, convertor.GetItemWithQuantityConvertor().EntitiesToProto(items))
	if err != nil {
		return nil, err
	}

	return convertor.GetItemConvertor().ProtoToEntities(resp.Items), nil
}

func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}
	var res []*entity.ItemWithQuantity
	for id, quantity := range merged {
		res = append(res, &entity.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}
	return res
}
