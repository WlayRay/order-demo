package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WlayRay/order-demo/common/broker"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/payment/app"
	"github.com/WlayRay/order-demo/payment/app/command"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
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
			c.handleMessage(ch, msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
	zap.L().Info("Received a new message", zap.String("queue", q.Name), zap.String("body", string(msg.Body)))

	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		zap.L().Warn("Failed to unmarshal message", zap.Error(err))
		_ = msg.Nack(false, false)
		return
	}

	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
		zap.L().Warn("Failed to handle message", zap.Error(err))
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			zap.L().Warn("Message retry error", zap.Error(err), zap.Any("messageID", msg.MessageId))
		}
		_ = msg.Nack(false, false)
		return
	}

	span.AddEvent("payment.created")
	_ = msg.Ack(false)
	zap.L().Info("consumed successfully")
}
