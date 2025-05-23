package adapters

import (
	"context"
	"fmt"
	"hash/fnv"
	"reflect"
	"strconv"
	"time"

	_ "github.com/WlayRay/order-demo/common/config"
	"github.com/WlayRay/order-demo/common/db"
	"github.com/WlayRay/order-demo/common/lib"
	domain "github.com/WlayRay/order-demo/order/domain/order"
	"github.com/WlayRay/order-demo/order/ent"
	orderModel "github.com/WlayRay/order-demo/order/ent/order"
	"github.com/WlayRay/order-demo/order/entity"
	_ "github.com/lib/pq" // 驱动导入
	"go.uber.org/zap"
)

type OrderRepositoryPG struct {
	db *ent.Client
}

func NewOrderRepositoryPG(db *ent.Client) *OrderRepositoryPG {
	return &OrderRepositoryPG{db: db}
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

// Create 创建订单并返回持久化后的订单对象
func (o OrderRepositoryPG) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// 将领域对象转换为存储层结构
	itemsData := make(map[string]any, len(order.Items))
	for i, item := range order.Items {
		itemsData[fmt.Sprintf("item-%d", i)] = &entity.Item{
			ID:       item.ID,
			Name:     item.Name,
			PriceID:  item.PriceID,
			Quantity: item.Quantity,
		}
	}

	ip, err := lib.GetLocalIP()
	if err != nil {
		return nil, fmt.Errorf("failed to get local IP: %w", err)
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(ip))
	snowflakeInstance, err := lib.GetSnowflakeInstance(h.Sum64()%1024, 10*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake snowflake instance: %w", err)
	}

	id, err := snowflakeInstance.GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get snowflake ID: %w", err)
	}

	created, err := o.db.Order.Create().
		SetOrderID(strconv.FormatUint(id, 10)).
		SetCustomerID(order.CustomerID).
		SetStatus(order.Status).
		SetPaymentLink(order.PaymentLink).
		SetItems(itemsData).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create order failed: %w", err)
	}

	// 返回带有数据库生成ID的领域对象
	return &domain.Order{
		ID:          created.OrderID, // 注意使用业务订单号
		CustomerID:  created.CustomerID,
		Status:      created.Status,
		PaymentLink: created.PaymentLink,
		Items:       order.Items, // 保持原始领域对象数据
	}, nil
}

// Get 根据业务订单ID和客户ID获取订单
func (o OrderRepositoryPG) Get(ctx context.Context, id, customerID string) (*domain.Order, error) {
	// 使用Ent查询订单
	entOrder, err := o.db.Order.Query().
		Where(
			orderModel.OrderID(id),
			orderModel.CustomerID(customerID),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	var items []*entity.Item
	for _, v := range entOrder.Items {
		if item, ok := v.(map[string]any); ok {
			items = append(items, &entity.Item{
				ID:       item["id"].(string),
				Name:     item["name"].(string),
				PriceID:  item["priceID"].(string),
				Quantity: int32(item["quantity"].(float64)),
			})
		} else {
			zap.L().Warn("unexpected item type", zap.Any("item", v), zap.Any("item type", reflect.TypeOf(v)))
		}
	}

	return &domain.Order{
		ID:          entOrder.OrderID,
		CustomerID:  entOrder.CustomerID,
		Status:      entOrder.Status,
		PaymentLink: entOrder.PaymentLink,
		Items:       items,
	}, nil
}

func (o OrderRepositoryPG) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
	// 开启事务
	tx, err := o.db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			zap.L().Error("更新订单时发生 panic", zap.Any("panic", p))
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 执行更新函数
	updatedOrder, err := updateFn(ctx, order)
	if err != nil {
		return fmt.Errorf("更新函数执行失败: %w", err)
	}

	// 转换商品数据
	itemsData := convertItemsToMap(updatedOrder)

	// 执行更新操作
	count, err := tx.Order.Update().
		Where(
			orderModel.OrderID(updatedOrder.ID),
			orderModel.CustomerID(updatedOrder.CustomerID),
		).
		SetStatus(updatedOrder.Status).
		SetPaymentLink(updatedOrder.PaymentLink).
		SetItems(itemsData).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("在事务中更新订单失败: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("更新订单失败：未找到匹配的记录，orderID: %s, customerID: %s", updatedOrder.ID, updatedOrder.CustomerID)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	zap.L().Info("订单更新成功",
		zap.String("orderID", updatedOrder.ID),
		zap.String("customerID", updatedOrder.CustomerID),
		zap.String("status", updatedOrder.Status))

	return nil
}

func convertItemsToMap(updatedOrder *domain.Order) map[string]any {
	itemsData := make(map[string]any, len(updatedOrder.Items))
	for i, item := range updatedOrder.Items {
		itemsData[fmt.Sprintf("item-%d", i)] = &entity.Item{
			ID:       item.ID,
			Name:     item.Name,
			PriceID:  item.PriceID,
			Quantity: item.Quantity,
		}
	}
	return itemsData
}
