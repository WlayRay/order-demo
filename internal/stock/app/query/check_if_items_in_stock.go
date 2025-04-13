package query

import (
	"context"
	"github.com/WlayRay/order-demo/common/decorator"
	domain "github.com/WlayRay/order-demo/stock/domain/stock"
	"github.com/WlayRay/order-demo/stock/entity"
	"github.com/WlayRay/order-demo/stock/infrastructure/integration"
	"go.uber.org/zap"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemInStockHandler struct {
	stockRepo domain.Repository
	stripeAPI *integration.StripeAPI
}

func NewCheckIfItemsInStockHandler(stockRepo domain.Repository,
	stripeAPI *integration.StripeAPI,
	logger *zap.Logger,
	metricClient decorator.MetricsClient) CheckIfItemsInStockHandler {
	if stripeAPI == nil {
		panic("stripeAPI is nil")
	}
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
		checkIfItemInStockHandler{stockRepo: stockRepo, stripeAPI: stripeAPI},
		logger,
		metricClient,
	)
}

//var priceIds = [3]string{
//	"price_1R7HVgPNegMNE0WfuwRkVr6b",
//	"price_1RD4V5PNegMNE0WfaN9nu9vo",
//	"price_1RD4XoPNegMNE0Wf9is4F4Wg",
//}

func (c checkIfItemInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
	var res []*entity.Item
	for i := range len(query.Items) {
		priceID, err := c.stripeAPI.GetPriceByProductID(ctx, query.Items[i].ID)
		if err != nil {
			zap.L().Warn("GetPriceByProductID", zap.String("productID", query.Items[i].ID), zap.Error(err))
			return nil, err
		}
		res = append(res, &entity.Item{
			ID:       query.Items[i].ID,
			Quantity: query.Items[i].Quantity,
			PriceID:  priceID,
		})

	}
	return res, nil
}
