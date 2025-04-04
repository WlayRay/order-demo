package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
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
	var items []*stripe.CheckoutSessionLineItemParams
	for i := range len(order.Items) {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String("price_1R7HVgPNegMNE0WfuwRkVr6b"),
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
