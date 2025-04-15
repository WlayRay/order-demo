package adapters

import (
	"context"
	"github.com/WlayRay/order-demo/common/db"
	"github.com/WlayRay/order-demo/stock/ent"
	stockModel "github.com/WlayRay/order-demo/stock/ent/stock"
	"github.com/WlayRay/order-demo/stock/entity"
)

type StockRepositoryPG struct {
	db *ent.Client
}

func NewStockRepositoryPG(db *ent.Client) *StockRepositoryPG {
	return &StockRepositoryPG{db: db}
}

func NewEntClient() *ent.Client {
	drv, err := db.GetPGSQLConn()
	if err != nil {
		panic(err)
	}

	client := ent.NewClient(
		ent.Driver(drv),
		ent.Debug(),
	)

	return client
}

func (s StockRepositoryPG) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (s StockRepositoryPG) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
	data, err := s.batchGetStockByID(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.ItemWithQuantity, 0, len(data))
	for _, item := range data {
		result = append(result, &entity.ItemWithQuantity{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		})
	}
	return result, nil
}

func (s StockRepositoryPG) batchGetStockByID(ctx context.Context, productIDs []string) ([]*ent.Stock, error) {
	return s.db.Stock.Query().
		Where(stockModel.ProductIDIn(productIDs...)).
		All(ctx)
}
