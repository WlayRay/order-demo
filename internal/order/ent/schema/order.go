package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	incrementalEnabled := true

	return []ent.Field{
		field.Int("id").
			Unique().
			Immutable().
			Annotations(
				entsql.Annotation{
					Incremental: &incrementalEnabled, // 启用自增
				}),

		field.String("order_id").
			SchemaType(map[string]string{
				dialect.Postgres: "VARCHAR(50)",
			}).
			Optional(),

		field.String("customer_id").
			SchemaType(map[string]string{
				dialect.Postgres: "VARCHAR(50)",
			}).
			NotEmpty(),

		field.String("status").
			SchemaType(map[string]string{
				dialect.Postgres: "VARCHAR(30)",
			}).
			Optional(),

		field.String("payment_link").
			SchemaType(map[string]string{
				dialect.Postgres: "VARCHAR(800)",
			}).
			Optional(),

		field.JSON("items", map[string]interface{}{}).
			SchemaType(map[string]string{
				dialect.Postgres: "JSONB",
			}),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return nil
}
