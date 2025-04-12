// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/WlayRay/order-demo/order/ent/order"
	"github.com/WlayRay/order-demo/order/ent/predicate"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeOrder = "Order"
)

// OrderMutation represents an operation that mutates the Order nodes in the graph.
type OrderMutation struct {
	config
	op            Op
	typ           string
	id            *int
	order_id      *string
	customer_id   *string
	status        *string
	payment_link  *string
	items         *map[string]interface{}
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Order, error)
	predicates    []predicate.Order
}

var _ ent.Mutation = (*OrderMutation)(nil)

// orderOption allows management of the mutation configuration using functional options.
type orderOption func(*OrderMutation)

// newOrderMutation creates new mutation for the Order entity.
func newOrderMutation(c config, op Op, opts ...orderOption) *OrderMutation {
	m := &OrderMutation{
		config:        c,
		op:            op,
		typ:           TypeOrder,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withOrderID sets the ID field of the mutation.
func withOrderID(id int) orderOption {
	return func(m *OrderMutation) {
		var (
			err   error
			once  sync.Once
			value *Order
		)
		m.oldValue = func(ctx context.Context) (*Order, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Order.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withOrder sets the old Order of the mutation.
func withOrder(node *Order) orderOption {
	return func(m *OrderMutation) {
		m.oldValue = func(context.Context) (*Order, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m OrderMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m OrderMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Order entities.
func (m *OrderMutation) SetID(id int) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *OrderMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *OrderMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Order.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetOrderID sets the "order_id" field.
func (m *OrderMutation) SetOrderID(s string) {
	m.order_id = &s
}

// OrderID returns the value of the "order_id" field in the mutation.
func (m *OrderMutation) OrderID() (r string, exists bool) {
	v := m.order_id
	if v == nil {
		return
	}
	return *v, true
}

// OldOrderID returns the old "order_id" field's value of the Order entity.
// If the Order object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OrderMutation) OldOrderID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldOrderID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldOrderID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldOrderID: %w", err)
	}
	return oldValue.OrderID, nil
}

// ClearOrderID clears the value of the "order_id" field.
func (m *OrderMutation) ClearOrderID() {
	m.order_id = nil
	m.clearedFields[order.FieldOrderID] = struct{}{}
}

// OrderIDCleared returns if the "order_id" field was cleared in this mutation.
func (m *OrderMutation) OrderIDCleared() bool {
	_, ok := m.clearedFields[order.FieldOrderID]
	return ok
}

// ResetOrderID resets all changes to the "order_id" field.
func (m *OrderMutation) ResetOrderID() {
	m.order_id = nil
	delete(m.clearedFields, order.FieldOrderID)
}

// SetCustomerID sets the "customer_id" field.
func (m *OrderMutation) SetCustomerID(s string) {
	m.customer_id = &s
}

// CustomerID returns the value of the "customer_id" field in the mutation.
func (m *OrderMutation) CustomerID() (r string, exists bool) {
	v := m.customer_id
	if v == nil {
		return
	}
	return *v, true
}

// OldCustomerID returns the old "customer_id" field's value of the Order entity.
// If the Order object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OrderMutation) OldCustomerID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCustomerID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCustomerID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCustomerID: %w", err)
	}
	return oldValue.CustomerID, nil
}

// ResetCustomerID resets all changes to the "customer_id" field.
func (m *OrderMutation) ResetCustomerID() {
	m.customer_id = nil
}

// SetStatus sets the "status" field.
func (m *OrderMutation) SetStatus(s string) {
	m.status = &s
}

// Status returns the value of the "status" field in the mutation.
func (m *OrderMutation) Status() (r string, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the Order entity.
// If the Order object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OrderMutation) OldStatus(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatus is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatus requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatus: %w", err)
	}
	return oldValue.Status, nil
}

// ClearStatus clears the value of the "status" field.
func (m *OrderMutation) ClearStatus() {
	m.status = nil
	m.clearedFields[order.FieldStatus] = struct{}{}
}

// StatusCleared returns if the "status" field was cleared in this mutation.
func (m *OrderMutation) StatusCleared() bool {
	_, ok := m.clearedFields[order.FieldStatus]
	return ok
}

// ResetStatus resets all changes to the "status" field.
func (m *OrderMutation) ResetStatus() {
	m.status = nil
	delete(m.clearedFields, order.FieldStatus)
}

// SetPaymentLink sets the "payment_link" field.
func (m *OrderMutation) SetPaymentLink(s string) {
	m.payment_link = &s
}

// PaymentLink returns the value of the "payment_link" field in the mutation.
func (m *OrderMutation) PaymentLink() (r string, exists bool) {
	v := m.payment_link
	if v == nil {
		return
	}
	return *v, true
}

// OldPaymentLink returns the old "payment_link" field's value of the Order entity.
// If the Order object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OrderMutation) OldPaymentLink(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPaymentLink is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPaymentLink requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPaymentLink: %w", err)
	}
	return oldValue.PaymentLink, nil
}

// ClearPaymentLink clears the value of the "payment_link" field.
func (m *OrderMutation) ClearPaymentLink() {
	m.payment_link = nil
	m.clearedFields[order.FieldPaymentLink] = struct{}{}
}

// PaymentLinkCleared returns if the "payment_link" field was cleared in this mutation.
func (m *OrderMutation) PaymentLinkCleared() bool {
	_, ok := m.clearedFields[order.FieldPaymentLink]
	return ok
}

// ResetPaymentLink resets all changes to the "payment_link" field.
func (m *OrderMutation) ResetPaymentLink() {
	m.payment_link = nil
	delete(m.clearedFields, order.FieldPaymentLink)
}

// SetItems sets the "items" field.
func (m *OrderMutation) SetItems(value map[string]interface{}) {
	m.items = &value
}

// Items returns the value of the "items" field in the mutation.
func (m *OrderMutation) Items() (r map[string]interface{}, exists bool) {
	v := m.items
	if v == nil {
		return
	}
	return *v, true
}

// OldItems returns the old "items" field's value of the Order entity.
// If the Order object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OrderMutation) OldItems(ctx context.Context) (v map[string]interface{}, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldItems is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldItems requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldItems: %w", err)
	}
	return oldValue.Items, nil
}

// ResetItems resets all changes to the "items" field.
func (m *OrderMutation) ResetItems() {
	m.items = nil
}

// Where appends a list predicates to the OrderMutation builder.
func (m *OrderMutation) Where(ps ...predicate.Order) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the OrderMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *OrderMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Order, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *OrderMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *OrderMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Order).
func (m *OrderMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *OrderMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.order_id != nil {
		fields = append(fields, order.FieldOrderID)
	}
	if m.customer_id != nil {
		fields = append(fields, order.FieldCustomerID)
	}
	if m.status != nil {
		fields = append(fields, order.FieldStatus)
	}
	if m.payment_link != nil {
		fields = append(fields, order.FieldPaymentLink)
	}
	if m.items != nil {
		fields = append(fields, order.FieldItems)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *OrderMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case order.FieldOrderID:
		return m.OrderID()
	case order.FieldCustomerID:
		return m.CustomerID()
	case order.FieldStatus:
		return m.Status()
	case order.FieldPaymentLink:
		return m.PaymentLink()
	case order.FieldItems:
		return m.Items()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *OrderMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case order.FieldOrderID:
		return m.OldOrderID(ctx)
	case order.FieldCustomerID:
		return m.OldCustomerID(ctx)
	case order.FieldStatus:
		return m.OldStatus(ctx)
	case order.FieldPaymentLink:
		return m.OldPaymentLink(ctx)
	case order.FieldItems:
		return m.OldItems(ctx)
	}
	return nil, fmt.Errorf("unknown Order field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *OrderMutation) SetField(name string, value ent.Value) error {
	switch name {
	case order.FieldOrderID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetOrderID(v)
		return nil
	case order.FieldCustomerID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCustomerID(v)
		return nil
	case order.FieldStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case order.FieldPaymentLink:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPaymentLink(v)
		return nil
	case order.FieldItems:
		v, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetItems(v)
		return nil
	}
	return fmt.Errorf("unknown Order field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *OrderMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *OrderMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *OrderMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Order numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *OrderMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(order.FieldOrderID) {
		fields = append(fields, order.FieldOrderID)
	}
	if m.FieldCleared(order.FieldStatus) {
		fields = append(fields, order.FieldStatus)
	}
	if m.FieldCleared(order.FieldPaymentLink) {
		fields = append(fields, order.FieldPaymentLink)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *OrderMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *OrderMutation) ClearField(name string) error {
	switch name {
	case order.FieldOrderID:
		m.ClearOrderID()
		return nil
	case order.FieldStatus:
		m.ClearStatus()
		return nil
	case order.FieldPaymentLink:
		m.ClearPaymentLink()
		return nil
	}
	return fmt.Errorf("unknown Order nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *OrderMutation) ResetField(name string) error {
	switch name {
	case order.FieldOrderID:
		m.ResetOrderID()
		return nil
	case order.FieldCustomerID:
		m.ResetCustomerID()
		return nil
	case order.FieldStatus:
		m.ResetStatus()
		return nil
	case order.FieldPaymentLink:
		m.ResetPaymentLink()
		return nil
	case order.FieldItems:
		m.ResetItems()
		return nil
	}
	return fmt.Errorf("unknown Order field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *OrderMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *OrderMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *OrderMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *OrderMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *OrderMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *OrderMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *OrderMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Order unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *OrderMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Order edge %s", name)
}
