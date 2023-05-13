package op

import "fmt"

func Or(cond1 *Condition, cond2 *Condition) *Condition {
	return &Condition{
		Query: fmt.Sprintf("(%s OR %s)", cond1.Query, cond2.Query),
		Value: []any{
			cond1.Value[0],
			cond2.Value[0],
		},
	}
}
