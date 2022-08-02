package repository_test

import (
	repository "github.com/jaztec/gorm-repository"
	"testing"
)

func TestNewCondition(t *testing.T) {
	t.Run("Create some conditions", func(t *testing.T) {
		type te struct {
			t repository.ConditionType
			q string
			a []any
		}
		tests := []te{
			{repository.TypeWhere, "a = ?", []any{"a"}},
			{repository.TypeJoin, "b = ? AND c = ? AND d = 1", []any{"a", 2}},
		}
		for _, ut := range tests {
			var c repository.Condition
			switch ut.t {
			case repository.TypeWhere:
				c = repository.NewWhereCondition(ut.q, ut.a)
			case repository.TypeJoin:
				c = repository.NewJoinCondition(ut.q, ut.a)
			default:
				t.Fatal("No valid type provided")
			}

			if c.Type() != ut.t {
				t.Errorf("Type %d does not equal %d", c.Type(), ut.t)
			}
			if c.Query() != ut.q {
				t.Errorf("Query %s does not equal %s", c.Query(), ut.q)
			}
			if len(c.Args()) != len(ut.a) {
				t.Errorf("Provided args have different count (%d) as condition (%d)", len(ut.a), len(c.Args()))
			}
			for i, _ := range ut.a {
				if c.Args()[i] != ut.a[i] {
					t.Error("Arguments at position do not match", i, ut.a[i], c.Args()[i])
				}
			}
		}
	})
}
