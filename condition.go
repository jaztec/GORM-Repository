package repository

type Condition struct {
	Query string
	Args  []any
}

func NewCondition(q string, a ...any) Condition {
	return Condition{
		q,
		a,
	}
}
