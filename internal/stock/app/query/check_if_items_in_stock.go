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
	if err := c.checkStock(ctx, query.Items); err != nil {
		return nil, err
	}

	var res []*entity.Item
	for i := range len(query.Items) {
		priceID, err := c.stripeAPI.GetPriceByProductID(ctx, query.Items[i].ID)
		if err != nil || priceID == "" {
			zap.L().Warn("GetPriceByProductID", zap.String("productID", query.Items[i].ID), zap.Error(err))
			return nil, err
		}
		res = append(res, &entity.Item{
			ID:       query.Items[i].ID,
			Quantity: query.Items[i].Quantity,
			PriceID:  priceID,
		})
	}
	// TODO: 扣减库存
	return res, nil
}

func (c checkIfItemInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
	ids := make([]string, 0, len(query))
	for _, item := range query {
		ids = append(ids, item.ID)
	}

	records, err := c.stockRepo.GetStock(ctx, ids)
	if err != nil {
		return err
	}

	idQuantityMap := make(map[string]int32, len(records))
	for _, record := range records {
		idQuantityMap[record.ID] += record.Quantity
	}

	var (
		ok        = true
		failedIDs []struct {
			ID   string
			Want int32
			Have int32
		}
	)

	for _, item := range query {
		if item.Quantity > idQuantityMap[item.ID] {
			ok = false
			failedIDs = append(failedIDs, struct {
				ID   string
				Want int32
				Have int32
			}{
				ID:   item.ID,
				Want: item.Quantity,
				Have: idQuantityMap[item.ID],
			})
			break
		}
	}
	if ok {
		return nil
	}
	return domain.ExceedStockError{FailedIDs: failedIDs}
}
