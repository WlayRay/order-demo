package adapters

import (
	"context"
	"fmt"
	"github.com/WlayRay/order-demo/common/db"
	"github.com/WlayRay/order-demo/stock/ent"
	stockModel "github.com/WlayRay/order-demo/stock/ent/stock"
	"github.com/WlayRay/order-demo/stock/entity"
	"go.uber.org/zap"
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

func (s StockRepositoryPG) UpdateStock(
	ctx context.Context,
	query []*entity.ItemWithQuantity,
	updateFunc func(context.Context, []*entity.ItemWithQuantity, []*entity.ItemWithQuantity) error,
) error {
	tx, err := s.db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 确保事务最终会被提交或回滚
	rollback := true
	defer func() {
		if rollback {
			if err := tx.Rollback(); err != nil {
				zap.L().Error("transaction rollback failed", zap.Error(err))
			}
		}
	}()

	currentStock, err := tx.Stock.Query().
		Where(stockModel.ProductIDIn(getProductIDs(query)...)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query stock failed: %w", err)
	}

	currentItems := make([]*entity.ItemWithQuantity, 0, len(currentStock))
	for _, s := range currentStock {
		currentItems = append(currentItems, &entity.ItemWithQuantity{
			ID:       s.ProductID,
			Quantity: s.Quantity,
		})
	}

	err = updateFunc(ctx, currentItems, query) // 更新currentItems
	if err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	for _, item := range currentItems {
		if err := tx.Stock.Update().
			Where(stockModel.ProductID(item.ID)).
			SetQuantity(item.Quantity).Exec(ctx); err != nil {
			return fmt.Errorf("update stock failed: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	rollback = false
	return nil
}

// 辅助函数：提取productIDs
func getProductIDs(items []*entity.ItemWithQuantity) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}
