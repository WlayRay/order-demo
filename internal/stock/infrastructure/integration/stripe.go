package integration

import (
	"context"
	_ "github.com/WlayRay/order-demo/common/config"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/product"
)

type StripeAPI struct {
	apikey string
}

func NewStripeAPI() *StripeAPI {
	return &StripeAPI{apikey: viper.GetString("stripe-key")}
}

func (s StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
	stripe.Key = s.apikey
	result, err := product.Get(pid, &stripe.ProductParams{})
	if err != nil {
		return "", err
	}
	return result.DefaultPrice.ID, err
}
