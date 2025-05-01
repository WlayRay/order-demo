package ports

import (
	"context"

	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/common/genproto/stockpb"
	"github.com/WlayRay/order-demo/common/tracing"
	"github.com/WlayRay/order-demo/stock/app" // 注意这里是stock
	"github.com/WlayRay/order-demo/stock/app/query"
	"github.com/WlayRay/order-demo/stock/entity"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	_, span := tracing.Start(ctx, "CheckIfItemsInStock")
	defer span.End()

	// TODO: 统一到convertor做转换
	entityItems := make([]*entity.ItemWithQuantity, 0, len(request.Items))
	for _, item := range request.Items {
		entityItems = append(entityItems, &entity.ItemWithQuantity{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}

	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: entityItems})
	if err != nil {
		return nil, err
	}

	// TODO: 统一到convertor做转换
	orderpbItems := make([]*orderpb.Item, 0, len(items))
	for _, item := range items {
		orderpbItems = append(orderpbItems, &orderpb.Item{
			ID:       item.ID,
			Name:     item.Name,
			PriceID:  item.PriceID,
			Quantity: item.Quantity,
		})
	}
	return &stockpb.CheckIfItemsInStockResponse{InStock: 1, Items: orderpbItems}, nil
}
