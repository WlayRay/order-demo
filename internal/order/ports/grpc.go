package ports

import (
	"context"

	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/order/app" // 注意这里是order
	"github.com/WlayRay/order-demo/order/app/command"
	"github.com/WlayRay/order-demo/order/app/query"
	"github.com/WlayRay/order-demo/order/convertor"
	domain "github.com/WlayRay/order-demo/order/domain/order"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	_, err := G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: request.CustomerID,
		Items:      convertor.GetItemWithQuantityConvertor().ProtoToEntities(request.Items),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	o, err := G.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerID: request.CustomerID,
		OrderID:    request.ID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertor.GetOrderConvertor().EntityToProto(o), err
}

func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
	zap.L().Info("UpdateOrder", zap.Any("request", request))
	order, newOrderErr := domain.NewOrder(
		request.ID,
		request.CustomerID,
		request.Status,
		request.PaymentLink,
		convertor.GetItemConvertor().ProtoToEntities(request.Items))
	if newOrderErr != nil {
		return nil, status.Error(codes.Internal, newOrderErr.Error())
	}

	_, err = G.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		},
	})
	return
}
