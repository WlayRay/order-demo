package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Stock holds the schema definition for the Stock entity.
type Stock struct {
	ent.Schema
}

// Fields of the Stock.
func (Stock) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("name").MaxLen(300).NotEmpty().Comment("名称"),
		field.String("price").MaxLen(50).NotEmpty().Comment("价格"),
		field.String("product_id").MaxLen(300).NotEmpty(),
		field.Int32("quantity").Min(0).Comment("库存"),
		field.Time("created_at").Default(func() time.Time { return time.Now() }),
		field.Time("updated_at").Default(func() time.Time { return time.Now() }).UpdateDefault(func() time.Time { return time.Now() }),
	}
}

// Edges of the Stock.
func (Stock) Edges() []ent.Edge {
	return nil
}
