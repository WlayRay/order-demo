package adapters

import (
	"context"
	domain "github.com/WlayRay/order-demo/order/domain/order"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

var fakeData = []*domain.Order{
	{
		ID:          "fake-ID",
		CustomerID:  "fake-customer-ID",
		Items:       nil,
		PaymentLink: "nil",
		Status:      "ok",
	},
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	s := make([]*domain.Order, 0, 1)
	s = append(s, fakeData...)
	return &MemoryOrderRepository{
		lock: &sync.RWMutex{},
		//store: make([]*domain.Order, 0, 100),
		store: s,
	}
}

func (m MemoryOrderRepository) Create(_ context.Context, _ *domain.Order) (*domain.Order, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	newOrder := &domain.Order{
		ID:          strconv.FormatInt(time.Now().UnixNano(), 10),
		CustomerID:  "",
		Items:       nil,
		PaymentLink: "",
		Status:      "",
	}
	m.store = append(m.store, newOrder)
	zap.L().Info("create order in memory", zap.Any("order:", newOrder))
	return newOrder, nil
}

func (m MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, o := range m.store {
		if o.ID == id && o.CustomerID == customerID {
			zap.L().Info("get order in memory", zap.Any("order:", o))
			return o, nil
		}
	}
	return nil, domain.NotFoundError{OrderID: id}
}

//func (m MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
//	m.lock.Lock()
//	defer m.lock.Unlock()
//
//	found := false
//	for i, o := range m.store {
//		if o.ID == order.ID && o.CustomerID == order.CustomerID {
//			found = true
//			updatedOrder, err := updateFn(ctx, o)
//			if err != nil {
//				return err
//			}
//			m.store[i] = updatedOrder
//		}
//	}
//
//	if !found {
//		return domain.NotFoundError{OrderID: order.ID}
//	}
//	return nil
//}

func (m MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	found := false
	for i, o := range m.store {
		if o.ID == order.ID && o.CustomerID == order.CustomerID {
			found = true
			updatedOrder, err := updateFn(ctx, o)
			if err != nil {
				return err
			}
			m.store[i] = updatedOrder
		}
	}

	if !found {
		return domain.NotFoundError{OrderID: order.ID}
	}
	return nil
}
