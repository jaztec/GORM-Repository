package repository_test

import (
	repository "github.com/jaztec/gorm-repository"
	"testing"
)

func TestNewCondition(t *testing.T) {
	t.Run("Create some conditions", func(t *testing.T) {
		type te struct {
			q string
			a []any
		}
		tests := []te{
			{"a = ?", []any{"a"}},
			{"b = ? AND c = ? AND d = 1", []any{"a", 2}},
		}
		for _, ut := range tests {
			c := repository.NewCondition(ut.q, ut.a...)

			if c.Query != ut.q {
				t.Errorf("Query %s does not equal %s", c.Query, ut.q)
			}
			if len(c.Args) != len(ut.a) {
				t.Errorf("Provided args have different count (%d) as condition (%d)", len(ut.a), len(c.Args))
			}
			for i, _ := range ut.a {
				if c.Args[i] != ut.a[i] {
					t.Error("Arguments at position do not match", i, ut.a[i], c.Args[i])
				}
			}
		}
	})
}
