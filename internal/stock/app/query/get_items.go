package query

import (
	"context"

	"github.com/WlayRay/order-demo/common/decorator"
	"github.com/WlayRay/order-demo/stock/domain"
	"github.com/WlayRay/order-demo/stock/entity"
	"go.uber.org/zap"
)

type GetItems struct {
	ItemIDs []string
}

type GetItemsHandler decorator.QueryHandler[GetItems, []*entity.Item]

type getItemsHandler struct {
	stockRepo domain.Repository
}

func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*entity.Item, error) {
	items, err := g.stockRepo.GetItems(ctx, query.ItemIDs)
	if err != nil {
		return nil, err
	}
	return items, err
}

func NewGetItemsHandler(stockRepo domain.Repository, logger *zap.Logger, metricClient decorator.MetricsClient) GetItemsHandler {
	if stockRepo != nil {
		return decorator.ApplyQueryDecorators(
			getItemsHandler{stockRepo: stockRepo},
			logger,
			metricClient,
		)
	}
	panic("stockRepo is nil")
}
