package broker

import (
	"context"
	"fmt"
	_ "github.com/WlayRay/order-demo/common/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"time"
)

const (
	DLX                = "dlx"
	DLQ                = "dlq"
	amqpRetryHeaderKey = "x-retry-header"
)

var maxRetryCount = viper.GetInt64("rabbitmq.max-retry-count")

// Connect establishes a connection to the RabbitMQ server and returns a channel.
func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, func() error { return nil }
	}

	ch, chErr := conn.Channel()
	if chErr != nil {
		zap.L().Fatal(chErr.Error())
	}

	err = ch.ExchangeDeclare(EventOrderCreated, "direct", true, false, false, false, nil)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	err = ch.ExchangeDeclare(EventOrderPaid, "fanout", true, false, false, false, nil)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	if err = createDLX(ch); err != nil {
		zap.L().Fatal(err.Error())
	}
	return ch, func() error { return conn.Close() }
}

func createDLX(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare("share_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(DLX, "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.QueueBind(q.Name, "", DLX, false, nil)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(DLQ, true, false, false, false, nil)
	return err
}

func HandleRetry(ctx context.Context, ch *amqp.Channel, d *amqp.Delivery) error {
	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}
	v, ok := d.Headers[amqpRetryHeaderKey]
	var retryCount int64
	if v != nil {
		retryCount, ok = v.(int64)
	} else {
		retryCount = 0
	}

	if !ok {
		retryCount = 0
	}
	retryCount++
	d.Headers[amqpRetryHeaderKey] = retryCount

	if retryCount > maxRetryCount {
		zap.L().Error("retry count exceeded, moving message to dlq", zap.Int("retryCount", int(retryCount)), zap.Any("message", d.MessageId))
		return ch.PublishWithContext(ctx, "", DLQ, false, false, amqp.Publishing{
			Headers:      d.Headers,
			ContentType:  "application/json",
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		})
	}
	time.Sleep(time.Duration(retryCount) * time.Second)
	return ch.PublishWithContext(ctx, d.Exchange, d.RoutingKey, false, false, amqp.Publishing{
		Headers:      d.Headers,
		ContentType:  "application/json",
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
	})
}

type RabbitMQHeaderCarrier map[string]any

func (r RabbitMQHeaderCarrier) Get(key string) string {
	value, ok := r[key]
	if !ok {
		return ""
	}
	return value.(string)
}

func (r RabbitMQHeaderCarrier) Set(key string, value string) {
	r[key] = value
}

func (r RabbitMQHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	return keys
}

func InjectRabbitMQHeaders(ctx context.Context) map[string]any {
	carrier := make(RabbitMQHeaderCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier
}

func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]any) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers))
}
