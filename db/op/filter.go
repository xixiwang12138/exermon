package op

type CondType uint8

type Condition struct {
	Query string
	Value []any
}

func binaryOp(field string, op string, value any) *Condition {
	return &Condition{Query: field + " " + op + " ?", Value: []any{value}}
}

// Eq 等于
func Eq(field string, value any) *Condition {
	return binaryOp(field, "=", value)
}

// Lt 小于
func Lt(field string, value any) *Condition {
	return binaryOp(field, "<", value)
}

// Le 小于等于
func Le(field string, value any) *Condition {
	return binaryOp(field, "<=", value)
}

// Gt 大于
func Gt(field string, value any) *Condition {
	return binaryOp(field, ">", value)
}

// Ge 大于等于
func Ge(field string, value any) *Condition {
	return binaryOp(field, ">=", value)
}

// Ne 不等于
func Ne(field string, value any) *Condition {
	return binaryOp(field, "!=", value)
}

// Like 模糊查询
func Like(field string, value string) *Condition {
	return binaryOp(field, "LIKE", value)
}

func In(field string, value any) *Condition {
	return binaryOp(field, "IN", value)
}

func NotIn(field string, value any) *Condition {
	return binaryOp(field, "NOT IN", value)
}
