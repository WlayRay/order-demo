package stock

import (
	"context"
	"fmt"
	"github.com/WlayRay/order-demo/stock/entity"
	"strings"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*entity.Item, error)
	GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error)
}

type NotFoundError struct {
	Missing []string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("not found items: %v", strings.Join(n.Missing, ","))
}

type ExceedStockError struct {
	FailedIDs []struct {
		ID   string
		Want int32
		Have int32
	}
}

func (e ExceedStockError) Error() string {
	var details []string
	for _, product := range e.FailedIDs {
		details = append(details, fmt.Sprintf("Product ID: %s, Want: %d, Have: %d", product.ID, product.Want, product.Have))
	}
	return fmt.Sprintf("these products do not have enough stock: {%s}", strings.Join(details, ", "))
}
