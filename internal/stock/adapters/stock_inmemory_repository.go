package adapters

import (
	"context"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	domain "github.com/WlayRay/order-demo/stock/domain/stock"
	//"go.uber.org/zap"
	"sync"
)

type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

var stub = map[string]*orderpb.Item{
	"item_id": {
		ID:       "123456",
		Name:     "袜子",
		PriceID:  "10086",
		Quantity: 100,
	},
}

func NewMemoryStockRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

func (m MemoryOrderRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var (
		res     []*orderpb.Item
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
