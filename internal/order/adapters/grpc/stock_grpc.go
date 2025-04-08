package grpc

import (
	"context"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/common/genproto/stockpb"
	"go.uber.org/zap"
)

type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}

}

func (s StockGRPC) CheckItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})
	zap.L().Info("stock grpc response", zap.Any("resp", resp))
	return resp, err
}

func (s StockGRPC) GetItems(ctx context.Context, itemsIDs []string) ([]*orderpb.Item, error) {
	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{ItemsIDs: itemsIDs})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}
