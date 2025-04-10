package convertor

import (
	"fmt"
	client "github.com/WlayRay/order-demo/common/client/order"
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	domain "github.com/WlayRay/order-demo/order/domain/order"
	"github.com/WlayRay/order-demo/order/entity"
)

type OrderConvertor struct {
}

type ItemConvertor struct {
}

type ItemWithQuantityConvertor struct {
}

func (c OrderConvertor) EntityToProto(o *domain.Order) *orderpb.Order {
	c.check(o)
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       GetItemConvertor().EntitiesToProto(o.Items),
	}
}

func (c OrderConvertor) ProtoToEntity(o *orderpb.Order) *domain.Order {
	c.check(o)
	return &domain.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       GetItemConvertor().ProtoToEntities(o.Items),
	}
}

func (c OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
	c.check(o)
	return &domain.Order{
		ID:          o.Id,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       GetItemConvertor().ClientsToEntities(o.Items),
	}
}

func (c OrderConvertor) EntityToClient(o *domain.Order) *client.Order {
	c.check(o)
	return &client.Order{
		Id:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       GetItemConvertor().EntitiesToClients(o.Items),
	}

}

func (c OrderConvertor) check(o any) {
	if o == nil {
		panic(fmt.Sprintf("convert failed, %T is nil", o))
	}
}

func (c ItemConvertor) EntitiesToProto(entity []*entity.Item) []*orderpb.Item {
	res := make([]*orderpb.Item, 0, len(entity))
	for _, item := range entity {
		res = append(res, c.EntityToProto(item))
	}
	return res
}

func (c ItemConvertor) ProtoToEntities(items []*orderpb.Item) []*entity.Item {
	res := make([]*entity.Item, 0, len(items))
	for _, item := range items {
		res = append(res, &entity.Item{
			ID:       item.ID,
			Name:     item.Name,
			PriceID:  item.PriceID,
			Quantity: item.Quantity,
		})
	}
	return res
}

func (c ItemConvertor) ClientsToEntities(items []client.Item) []*entity.Item {
	res := make([]*entity.Item, 0, len(items))
	for _, item := range items {
		res = append(res, &entity.Item{
			ID:       item.Id,
			Name:     item.Name,
			PriceID:  item.PriceID,
			Quantity: item.Quantity,
		})
	}
	return res
}

func (c ItemConvertor) EntitiesToClients(items []*entity.Item) []client.Item {
	res := make([]client.Item, 0, len(items))
	for _, item := range items {
		res = append(res, client.Item{
			Id:       item.ID,
			Name:     item.Name,
			PriceID:  item.PriceID,
			Quantity: item.Quantity,
		})
	}
	return res
}

func (c ItemConvertor) EntityToProto(i *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       i.ID,
		Name:     i.Name,
		PriceID:  i.PriceID,
		Quantity: i.Quantity,
	}
}

func (c ItemWithQuantityConvertor) EntitiesToProto(items []*entity.ItemWithQuantity) []*orderpb.ItemWithQuantity {
	res := make([]*orderpb.ItemWithQuantity, 0, len(items))
	for _, item := range items {
		res = append(res, c.EntityToProto(item))
	}
	return res
}

func (c ItemWithQuantityConvertor) EntityToProto(items *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
	return &orderpb.ItemWithQuantity{
		ID:       items.ID,
		Quantity: items.Quantity,
	}
}

func (c ItemWithQuantityConvertor) ProtoToEntities(items []*orderpb.ItemWithQuantity) []*entity.ItemWithQuantity {
	res := make([]*entity.ItemWithQuantity, 0, len(items))
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return res
}

func (c ItemWithQuantityConvertor) ProtoToEntity(items *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       items.ID,
		Quantity: items.Quantity,
	}
}

func (c ItemWithQuantityConvertor) ClientsToEntities(items []client.ItemWithQuantity) []*entity.ItemWithQuantity {
	res := make([]*entity.ItemWithQuantity, 0, len(items))
	for _, item := range items {
		res = append(res, c.ClientToEntity(item))
	}
	return res
}

func (c ItemWithQuantityConvertor) ClientToEntity(items client.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       items.Id,
		Quantity: items.Quantity,
	}
}
