package op

import "strings"

func boolOp(op string, cond1 *Condition, cond2 *Condition) *Condition {
	return &Condition{
		Query: strings.Join([]string{"(", cond1.Query, ") ", op, " (", cond2.Query, ")"}, ""),
		Value: append(cond1.Value, cond2.Value...),
	}
}

func And(cond1 *Condition, cond2 *Condition) *Condition {
	return boolOp("AND", cond1, cond2)
}

func Or(cond1 *Condition, cond2 *Condition) *Condition {
	return boolOp("OR", cond1, cond2)
}
