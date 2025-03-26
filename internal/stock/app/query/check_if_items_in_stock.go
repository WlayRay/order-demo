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
	var res []*orderpb.Item
	for i := 0; i < len(query.Items); i++ {
		res = append(
			res,
			&orderpb.Item{
				ID:       query.Items[i].ID,
				Quantity: query.Items[i].Quantity,
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
