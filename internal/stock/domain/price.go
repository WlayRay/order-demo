package domain

import "context"

type StripeAPIInterface interface {
	GetPriceByProductID(ctx context.Context, productID string) (string, error)
}
