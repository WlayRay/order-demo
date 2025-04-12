package query

import (
	"context"
	"github.com/WlayRay/order-demo/common/decorator"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	domain "github.com/WlayRay/order-demo/stock/domain/stock"
	"go.uber.org/zap"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemInStockHandler struct {
	stockRepo domain.Repository
}

func (c checkIfItemInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var priceIds = [3]string{
		"price_1R7HVgPNegMNE0WfuwRkVr6b",
		"price_1RD4V5PNegMNE0WfaN9nu9vo",
		"price_1RD4XoPNegMNE0Wf9is4F4Wg",
	}

	var res []*orderpb.Item
	for i := 0; i < len(query.Items); i++ {
		res = append(
			res,
			&orderpb.Item{
				ID:       query.Items[i].ID,
				Quantity: query.Items[i].Quantity,
				PriceID:  priceIds[i],
			})

	}
	return res, nil
}

func NewCheckIfItemsInStockHandler(stockRepo domain.Repository, logger *zap.Logger, metricClient decorator.MetricsClient) CheckIfItemsInStockHandler {
	if stockRepo != nil {
		return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
			checkIfItemInStockHandler{stockRepo: stockRepo},
			logger,
			metricClient,
		)
	}
	panic("stockRepo is nil")
}
