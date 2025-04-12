// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/WlayRay/order-demo/order/ent/order"
)

// OrderCreate is the builder for creating a Order entity.
type OrderCreate struct {
	config
	mutation *OrderMutation
	hooks    []Hook
}

// SetOrderID sets the "order_id" field.
func (oc *OrderCreate) SetOrderID(s string) *OrderCreate {
	oc.mutation.SetOrderID(s)
	return oc
}

// SetNillableOrderID sets the "order_id" field if the given value is not nil.
func (oc *OrderCreate) SetNillableOrderID(s *string) *OrderCreate {
	if s != nil {
		oc.SetOrderID(*s)
	}
	return oc
}

// SetCustomerID sets the "customer_id" field.
func (oc *OrderCreate) SetCustomerID(s string) *OrderCreate {
	oc.mutation.SetCustomerID(s)
	return oc
}

// SetStatus sets the "status" field.
func (oc *OrderCreate) SetStatus(s string) *OrderCreate {
	oc.mutation.SetStatus(s)
	return oc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (oc *OrderCreate) SetNillableStatus(s *string) *OrderCreate {
	if s != nil {
		oc.SetStatus(*s)
	}
	return oc
}

// SetPaymentLink sets the "payment_link" field.
func (oc *OrderCreate) SetPaymentLink(s string) *OrderCreate {
	oc.mutation.SetPaymentLink(s)
	return oc
}

// SetNillablePaymentLink sets the "payment_link" field if the given value is not nil.
func (oc *OrderCreate) SetNillablePaymentLink(s *string) *OrderCreate {
	if s != nil {
		oc.SetPaymentLink(*s)
	}
	return oc
}

// SetItems sets the "items" field.
func (oc *OrderCreate) SetItems(m map[string]interface{}) *OrderCreate {
	oc.mutation.SetItems(m)
	return oc
}

// SetID sets the "id" field.
func (oc *OrderCreate) SetID(i int) *OrderCreate {
	oc.mutation.SetID(i)
	return oc
}

// Mutation returns the OrderMutation object of the builder.
func (oc *OrderCreate) Mutation() *OrderMutation {
	return oc.mutation
}

// Save creates the Order in the database.
func (oc *OrderCreate) Save(ctx context.Context) (*Order, error) {
	return withHooks(ctx, oc.sqlSave, oc.mutation, oc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (oc *OrderCreate) SaveX(ctx context.Context) *Order {
	v, err := oc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (oc *OrderCreate) Exec(ctx context.Context) error {
	_, err := oc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (oc *OrderCreate) ExecX(ctx context.Context) {
	if err := oc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (oc *OrderCreate) check() error {
	if _, ok := oc.mutation.CustomerID(); !ok {
		return &ValidationError{Name: "customer_id", err: errors.New(`ent: missing required field "Order.customer_id"`)}
	}
	if v, ok := oc.mutation.CustomerID(); ok {
		if err := order.CustomerIDValidator(v); err != nil {
			return &ValidationError{Name: "customer_id", err: fmt.Errorf(`ent: validator failed for field "Order.customer_id": %w`, err)}
		}
	}
	if _, ok := oc.mutation.Items(); !ok {
		return &ValidationError{Name: "items", err: errors.New(`ent: missing required field "Order.items"`)}
	}
	return nil
}

func (oc *OrderCreate) sqlSave(ctx context.Context) (*Order, error) {
	if err := oc.check(); err != nil {
		return nil, err
	}
	_node, _spec := oc.createSpec()
	if err := sqlgraph.CreateNode(ctx, oc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int(id)
	}
	oc.mutation.id = &_node.ID
	oc.mutation.done = true
	return _node, nil
}

func (oc *OrderCreate) createSpec() (*Order, *sqlgraph.CreateSpec) {
	var (
		_node = &Order{config: oc.config}
		_spec = sqlgraph.NewCreateSpec(order.Table, sqlgraph.NewFieldSpec(order.FieldID, field.TypeInt))
	)
	if id, ok := oc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := oc.mutation.OrderID(); ok {
		_spec.SetField(order.FieldOrderID, field.TypeString, value)
		_node.OrderID = value
	}
	if value, ok := oc.mutation.CustomerID(); ok {
		_spec.SetField(order.FieldCustomerID, field.TypeString, value)
		_node.CustomerID = value
	}
	if value, ok := oc.mutation.Status(); ok {
		_spec.SetField(order.FieldStatus, field.TypeString, value)
		_node.Status = value
	}
	if value, ok := oc.mutation.PaymentLink(); ok {
		_spec.SetField(order.FieldPaymentLink, field.TypeString, value)
		_node.PaymentLink = value
	}
	if value, ok := oc.mutation.Items(); ok {
		_spec.SetField(order.FieldItems, field.TypeJSON, value)
		_node.Items = value
	}
	return _node, _spec
}

// OrderCreateBulk is the builder for creating many Order entities in bulk.
type OrderCreateBulk struct {
	config
	err      error
	builders []*OrderCreate
}

// Save creates the Order entities in the database.
func (ocb *OrderCreateBulk) Save(ctx context.Context) ([]*Order, error) {
	if ocb.err != nil {
		return nil, ocb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ocb.builders))
	nodes := make([]*Order, len(ocb.builders))
	mutators := make([]Mutator, len(ocb.builders))
	for i := range ocb.builders {
		func(i int, root context.Context) {
			builder := ocb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*OrderMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ocb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ocb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ocb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ocb *OrderCreateBulk) SaveX(ctx context.Context) []*Order {
	v, err := ocb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ocb *OrderCreateBulk) Exec(ctx context.Context) error {
	_, err := ocb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ocb *OrderCreateBulk) ExecX(ctx context.Context) {
	if err := ocb.Exec(ctx); err != nil {
		panic(err)
	}
}
