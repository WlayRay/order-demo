package ports

import (
	"context"

	"github.com/WlayRay/order-demo/common/genproto/stockpb"
	"github.com/WlayRay/order-demo/common/tracing"
	"github.com/WlayRay/order-demo/stock/app" // 注意这里是stock
	"github.com/WlayRay/order-demo/stock/app/query"
	"github.com/WlayRay/order-demo/stock/convertor"
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

	entityItems := convertor.GetItemConvertor().ProtoToEntities(request.Items)
	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: entityItems})
	if err != nil {
		return nil, err
	}

	orderpbItems := convertor.GetOrderConvertor().EntitiesToProto(items)
	return &stockpb.CheckIfItemsInStockResponse{InStock: 1, Items: orderpbItems}, nil
}
