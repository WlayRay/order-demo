package query

import (
	"context"
	"fmt"
	"github.com/WlayRay/order-demo/common/db"
	"github.com/WlayRay/order-demo/common/decorator"
	domain "github.com/WlayRay/order-demo/stock/domain/stock"
	"github.com/WlayRay/order-demo/stock/entity"
	"github.com/WlayRay/order-demo/stock/infrastructure/integration"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

const ETCDLockPrefix = "/stock/lock/"

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
	session, mutex, err := lock(ctx, getLockKey(query))
	if err != nil {
		return nil, err
	}
	defer func() {
		var releaseErr error
		releaseErr = mutex.Unlock(ctx)
		releaseErr = session.Close()
		if releaseErr != nil {
			zap.L().Warn("etcd unlock failed", zap.Error(releaseErr))
		}
	}()

	res := make([]*entity.Item, 0, len(query.Items))
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
	// TODO: 拆分出扣减库存的逻辑（如果需要的话）
	if err := c.checkStock(ctx, query.Items); err != nil {
		return nil, err
	}
	return res, nil
}

var lockKeyBuilderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

func getLockKey(query CheckIfItemsInStock) string {
	builder := lockKeyBuilderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		lockKeyBuilderPool.Put(builder)
	}()

	builder.WriteString(ETCDLockPrefix)
	for _, item := range query.Items {
		builder.WriteByte('-')
		builder.WriteString(item.ID)
	}
	return builder.String()
}

func lock(ctx context.Context, key string) (*concurrency.Session, *concurrency.Mutex, error) {
	etcdClient, _ := db.GetEtcdClient()
	timeoutCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if session, err := concurrency.NewSession(etcdClient); err != nil {
		return nil, nil, err
	} else {
		mutex := concurrency.NewMutex(session, key)
		if err := mutex.Lock(timeoutCtx); err != nil {
			return nil, nil, err
		} else {
			return session, mutex, nil
		}
	}
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
		return c.stockRepo.UpdateStock(ctx, query, func(
			ctx context.Context,
			existing []*entity.ItemWithQuantity,
			query []*entity.ItemWithQuantity,
		) error {
			// 创建现有库存的映射，提高查找效率
			stockMap := make(map[string]*entity.ItemWithQuantity, len(existing))
			for _, item := range existing {
				stockMap[item.ID] = item
			}

			for _, item := range query {
				existingItem, ok := stockMap[item.ID]
				if !ok {
					return fmt.Errorf("商品 %s 不存在", item.ID)
				}

				// 再次验证库存是否充足
				if existingItem.Quantity < item.Quantity {
					return fmt.Errorf("商品 %s 库存不足，当前库存: %d, 需求数量: %d",
						item.ID, existingItem.Quantity, item.Quantity)
				}
				existingItem.Quantity -= item.Quantity
			}

			return nil
		})
	}

	return domain.ExceedStockError{FailedIDs: failedIDs}
}
