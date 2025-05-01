package consumer

import "C"
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/WlayRay/order-demo/common/broker"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, request *orderpb.Order) error
}

type Consumer struct {
	orderGRPC OrderService
}

type order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}

func NewConsumer(orderGRPC OrderService) *Consumer {
	return &Consumer{
		orderGRPC: orderGRPC,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		zap.L().Fatal("Failed to declare queue", zap.Error(err))
	}

	if err := ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil); err != nil {
		zap.L().Fatal("Failed to bind queue", zap.Error(err))
	}

	msgs, consumeErr := ch.Consume(q.Name, "", false, false, false, false, nil)
	if consumeErr != nil {
		zap.L().Warn("Failed to consume messages", zap.Error(consumeErr))
	}

	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			c.handleMessage(ch, msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
	var err error
	zap.L().Info("Received a new message", zap.String("queue", q.Name), zap.String("body", string(msg.Body)))
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	mqCtx, span := tr.Start(context.Background(), fmt.Sprintf("rabbitmq.%s.consume", q.Name))

	defer func() {
		span.End()
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()

	o := &order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		zap.L().Warn("Failed to unmarshal message", zap.Error(err))
	}
	if o.Status != "paid" {
		err = errors.New("order has not paid yet")
	}
	outbound(o)
	span.AddEvent(fmt.Sprintf("order.cook: %+v", o))
	if err := c.orderGRPC.UpdateOrder(ctx, &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      "ready",
		PaymentLink: o.PaymentLink,
		Items:       o.Items,
	}); err != nil {
		if err = broker.HandleRetry(mqCtx, ch, &msg); err != nil {
			zap.L().Warn("Goods: failed to handle retry", zap.Error(err))
		}
		return
	}

	span.AddEvent("goods.order.finished.updated")
	zap.L().Info("consumed successfully")
}

func outbound(o *order) {
	zap.L().Debug("order id paid, goods ready to outbound", zap.Int("goods num", len(o.Items)))
	time.Sleep(3 * time.Second)
	zap.L().Debug("outbound success", zap.String("order id", o.ID), zap.String("customer id", o.CustomerID))
}
