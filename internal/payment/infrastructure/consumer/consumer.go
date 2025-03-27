package consumer

import (
	"context"
	"encoding/json"
	"github.com/WlayRay/order-demo/common/broker"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/payment/app"
	"github.com/WlayRay/order-demo/payment/app/command"
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
	q, queueErr := ch.QueueDeclare(broker.EventOrderCreated, false, true, false, false, nil)
	if queueErr != nil {
		zap.L().Fatal("Failed to declare a queue", zap.Error(queueErr))
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		zap.L().Warn("Failed to consume", zap.String("queue", q.Name), zap.Error(err))
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
	zap.L().Info("Received a new message", zap.String("queue", q.Name), zap.String("body", string(msg.Body)))

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		zap.L().Warn("Failed to unmarshal message", zap.Error(err))
		_ = msg.Nack(false, false)
		return
	}

	if _, err := c.app.Commands.CreatePaymentLink.Handle(context.TODO(), command.CreatePaymentLink{Order: o}); err != nil {
		//TODO: retry
		zap.L().Warn("Failed to handle message", zap.Error(err))
		_ = msg.Nack(false, false)
		return
	}

	_ = msg.Ack(false)
	zap.L().Info("consumed successfully")
}
