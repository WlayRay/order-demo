package order

import (
	"errors"
	"fmt"
	"github.com/WlayRay/order-demo/order/entity"
	"github.com/stripe/stripe-go/v80"
)

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*entity.Item
}

func NewOrder(id, customerID, status, paymentLink string, items []*entity.Item) (*Order, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	if customerID == "" {
		return nil, errors.New("empty customerID")
	}
	if items == nil {
		return nil, errors.New("empty items")
	}
	return &Order{
		ID:          id,
		CustomerID:  customerID,
		Status:      status,
		PaymentLink: paymentLink,
		Items:       items,
	}, nil
}

func (o Order) IsPaid() error {
	if o.Status != string(stripe.CheckoutSessionPaymentStatusPaid) {
		return fmt.Errorf("order %s is not paid, status = %s", o.ID, o.Status)
	}
	return nil
}
