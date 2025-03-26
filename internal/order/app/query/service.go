package query

import (
	"context"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/common/genproto/stockpb"
)

type StockService interface {
	CheckItemsInStock(items []*orderpb.ItemWithQuantity, ctx context.Context) (*stockpb.CheckIfItemsInStockResponse, error)
	GetItems(ctx context.Context, itemsIDs []string) ([]*orderpb.Item, error)
}
