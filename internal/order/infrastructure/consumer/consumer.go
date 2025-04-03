package consumer

import (
	"context"
	"encoding/json"
	"github.com/WlayRay/order-demo/common/broker"
	"github.com/WlayRay/order-demo/order/app"
	"github.com/WlayRay/order-demo/order/app/command"
	domain "github.com/WlayRay/order-demo/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{
		app: app,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderPaid, true, true, true, false, nil)
	if err != nil {
		zap.L().Fatal("Failed to declare a queue", zap.Error(err))
	}

	err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil)
	if err != nil {
		zap.L().Fatal("Failed to bind a queue", zap.Error(err))
	}
	
	msgs, consumeErr := ch.Consume(q.Name, "", false, false, false, false, nil)
	if consumeErr != nil {
		zap.L().Fatal("Failed to consume", zap.String("queue", q.Name), zap.Error(consumeErr))
	}

	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q, ch)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
	o := &domain.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		zap.L().Warn("Failed to unmarshal message", zap.Error(err))
		_ = msg.Nack(false, false)
		return
	}

	_, err := c.app.Commands.UpdateOrder.Handle(context.TODO(), command.UpdateOrder{
		Order: o,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.IsPaid(); err != nil {
				return nil, err
			}
			return order, nil
		},
	})
	if err != nil {
		zap.L().Warn("Failed to update order", zap.Error(err), zap.String("orderID", o.ID))
		// TODO: retry
		_ = msg.Nack(false, false)
		return
	}

	_ = msg.Ack(false)
	zap.L().Info("Order consume paid event success", zap.String("orderID", o.ID))
}
