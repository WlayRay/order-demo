package domain

import "context"

type PaymentInterface interface {
	GetPriceByProductID(ctx context.Context, productID string) (string, error)
}
