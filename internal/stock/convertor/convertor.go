package convertor

import (
	"github.com/WlayRay/order-demo/common/genproto/orderpb"
	"github.com/WlayRay/order-demo/stock/entity"
)

type ItemConvertor struct {
}

type OrderConvertor struct {
}

func (ic *ItemConvertor) ProtoToEntities(items []*orderpb.ItemWithQuantity) []*entity.ItemWithQuantity {
	result := make([]*entity.ItemWithQuantity, len(items))
	for i, item := range items {
		result[i] = ic.ProtoToEntity(item)
	}
	return result
}

func (*ItemConvertor) ProtoToEntity(o *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       o.ID,
		Quantity: o.Quantity,
	}
}

func (oc *OrderConvertor) EntitiesToProto(entities []*entity.Item) []*orderpb.Item {
	result := make([]*orderpb.Item, len(entities))
	for i, item := range entities {
		result[i] = oc.EntityToProto(item)
	}
	return result
}

func (*OrderConvertor) EntityToProto(o *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       o.ID,
		Name:     o.Name,
		Quantity: o.Quantity,
		PriceID:  o.PriceID,
	}
}
