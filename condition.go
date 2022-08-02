package repository

type ConditionType int

const (
	TypeWhere = iota
	TypeJoin
)

type Condition interface {
	Type() ConditionType
	Query() string
	Args() []any
}

type condition struct {
	t ConditionType
	q string
	a []any
}

func (c condition) Type() ConditionType {
	return c.t
}

func (c condition) Query() string {
	return c.q
}

func (c condition) Args() []any {
	return c.a
}

func NewWhereCondition(q string, a ...any) Condition {
	return condition{
		t: TypeWhere,
		q: q,
		a: a,
	}
}

func NewJoinCondition(q string, a ...any) Condition {
	return condition{
		t: TypeWhere,
		q: q,
		a: a,
	}
}
