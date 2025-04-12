package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/common/tracing"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/checkout/session"
)

type StripeProcessor struct {
	apikey string
}

func NewStripeProcessor(apikey string) *StripeProcessor {
	if apikey == "" {
		panic("empty api key")
	}
	stripe.Key = apikey
	return &StripeProcessor{apikey: apikey}
}

const (
	successURL = "http://localhost:9000/success"
	cancelURL  = "http://localhost:9000/cancel"
)

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	_, span := tracing.Start(ctx, "stripeProcessor.CreatePaymentLink")
	defer span.End()

	var items []*stripe.CheckoutSessionLineItemParams
	var priceIds = [3]string{
		"price_1R7HVgPNegMNE0WfuwRkVr6b",
		"price_1RD4V5PNegMNE0WfaN9nu9vo",
		"price_1RD4XoPNegMNE0Wf9is4F4Wg",
	}
	for i := range len(order.Items) {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(priceIds[i]),
			Quantity: stripe.Int64(int64(order.Items[i].Quantity)),
		})
	}

	marshalItems, _ := json.Marshal(order.Items)
	metadata := map[string]string{
		"orderID":     order.ID,
		"customerID":  order.CustomerID,
		"status":      order.Status,
		"items":       string(marshalItems),
		"paymentLink": order.PaymentLink,
	}
	params := &stripe.CheckoutSessionParams{
		Metadata:   metadata,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s?customerID=%s&orderID=%s", successURL, order.CustomerID, order.ID)),
		CancelURL:  stripe.String(cancelURL),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}
