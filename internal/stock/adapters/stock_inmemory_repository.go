package adapters

import (
	"context"
	domain "github.com/WlayRay/order-demo/stock/domain/stock"
	"github.com/WlayRay/order-demo/stock/entity"

	//"go.uber.org/zap"
	"sync"
)

type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*entity.Item
}

var stub = map[string]*entity.Item{
	"item1": {
		ID:       "123456",
		Name:     "袜子",
		PriceID:  "10086",
		Quantity: 100,
	},
	"item2": {
		ID:       "123456",
		Name:     "裤子",
		PriceID:  "10086",
		Quantity: 200,
	},
	"item3": {
		ID:       "123456",
		Name:     "机子",
		PriceID:  "10086",
		Quantity: 300,
	},
}

func NewMemoryStockRepository() *MemoryStockRepository {
	return &MemoryStockRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var (
		res     []*entity.Item
		missing []string
	)
	for i := 0; i < len(ids); i++ {
		if item, exist := m.store[ids[i]]; exist {
			res = append(res, item)
		} else {
			missing = append(missing, ids[i])
		}
	}
	if len(res) != len(ids) {
		return nil, domain.NotFoundError{Missing: missing}
	}
	return res, nil
}
