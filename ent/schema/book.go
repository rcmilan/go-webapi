package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

type Book struct{ ent.Schema }

func (Book) Fields() []ent.Field {
	return []ent.Field{
		field.Uint32("id").
			SchemaType(map[string]string{dialect.MySQL: "int unsigned"}).
			Positive().
			Immutable(),
		field.String("title").MaxLen(255).NotEmpty(),
		field.String("isbn").MaxLen(13).Unique().NotEmpty(),
		field.Float("price").Positive(),
		field.Int64("created_at").
			Immutable().
			DefaultFunc(func() int64 { return time.Now().Unix() }),
		field.Int("release_year").Positive(),
	}
}
