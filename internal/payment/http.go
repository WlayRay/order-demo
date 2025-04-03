package main

import (
	"context"
	"encoding/json"
	"github.com/WlayRay/order-demo/common/broker"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/payment/domain"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/webhook"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type PaymentHandler struct {
	channel *amqp.Channel
}

func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
	return &PaymentHandler{channel: ch}
}

func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", h.handleWebhook)
}

func (h *PaymentHandler) handleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zap.L().Info("Failed to read request body: %v", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, err)
		return
	}

	event, constructErr := webhook.ConstructEvent(payload,
		c.Request.Header.Get("Stripe-Signature"),
		viper.GetString("ENDPOINT_STRIPE_SECRET"))
	zap.L().Info("",
		zap.Any("Stripe-Signature", c.Request.Header.Get("Stripe-Signature")),
		zap.String("endpoint_stripe_secret", viper.GetString("ENDPOINT_STRIPE_SECRET")))

	if constructErr != nil {
		zap.L().Info("Failed to construct webhook event: %v", zap.Error(constructErr))
		c.JSON(http.StatusServiceUnavailable, constructErr)
		return
	}

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			zap.L().Info("Failed to unmarshal event data: %v", zap.Error(err))
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			zap.L().Info("Payment status is paid", zap.String("session", session.ID))

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var items []*orderpb.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)
			marshalledOrder, err := json.Marshal(&domain.Order{
				ID:          session.Metadata["orderID"],
				CustomerID:  session.Metadata["customerID"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
				PaymentLink: session.Metadata["paymentLink"],
				Items:       items,
			})
			if err != nil {
				zap.L().Info("Failed to marshal order: %v", zap.Error(err))
				c.JSON(http.StatusInternalServerError, err)
				return
			}

			err = h.channel.PublishWithContext(ctx,
				broker.EventOrderPaid,
				broker.EventOrderPaid,
				false,
				false,
				amqp.Publishing{
					ContentType:  "application/json",
					DeliveryMode: amqp.Persistent,
					Body:         marshalledOrder,
				})
			if err != nil {
				zap.L().Error("Failed to publish message", zap.Error(err))
				c.JSON(http.StatusInternalServerError, err)
				return
			}
		}
		c.JSON(http.StatusOK, nil)
	}
}
