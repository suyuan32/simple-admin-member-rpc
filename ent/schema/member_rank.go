package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/suyuan32/simple-admin-common/orm/ent/mixins"
)

type MemberRank struct {
	ent.Schema
}

func (MemberRank) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("Rank name | 等级名称"),
		field.String("code").Comment("Rank code | 等级码"),
		field.String("description").Comment("Rank description | 等级描述"),
		field.String("remark").Comment("Remark | 备注"),
	}
}

func (MemberRank) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.IDMixin{},
	}
}

func (MemberRank) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("members", Member.Type).Ref("ranks"),
	}
}

func (MemberRank) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").Unique(),
	}
}

func (MemberRank) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "mms_ranks"},
	}
}
