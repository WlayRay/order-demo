package query

import (
	"context"
	"github.com/WlayRay/order-demo/common/decorator"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	domain "github.com/WlayRay/order-demo/stock/domain/stock"
	"go.uber.org/zap"
)

type GetItems struct {
	ItemIDs []string
}

type GetItemsHandler decorator.QueryHandler[GetItems, []*orderpb.Item]

type getItemsHandler struct {
	stockRepo domain.Repository
}

func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*orderpb.Item, error) {
	items, err := g.stockRepo.GetItems(ctx, query.ItemIDs)
	if err != nil {
		return nil, err
	}
	return items, err
}

func NewGetItemsHandler(stockRepo domain.Repository, logger *zap.Logger, metricClient decorator.MetricsClient) GetItemsHandler {
	if stockRepo != nil {
		return decorator.ApplyQueryDecorators[GetItems, []*orderpb.Item](
			getItemsHandler{stockRepo: stockRepo},
			logger,
			metricClient,
		)
	}
	panic("stockRepo is nil")
}
