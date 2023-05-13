package op

type CondType uint8

type Condition struct {
	Query string
	Value []any
}

// Eq 等于
func Eq(field string, value any) *Condition {
	return &Condition{Query: field + " = ?", Value: []any{value}}
}

// Lt 小于
func Lt(field string, value any) *Condition {
	return &Condition{Query: field + " < ?", Value: []any{value}}
}

// Le 小于等于
func Le(field string, value any) *Condition {
	return &Condition{Query: field + " <= ?", Value: []any{value}}
}

// Gt 大于
func Gt(field string, value any) *Condition {
	return &Condition{Query: field + " > ?", Value: []any{value}}
}

// Ge 大于等于
func Ge(field string, value any) *Condition {
	return &Condition{Query: field + " >= ?", Value: []any{value}}
}

// Ne 不等于
func Ne(field string, value any) *Condition {
	return &Condition{Query: field + " <> ?", Value: []any{value}}
}

// Like 模糊查询
func Like(field string, value string) *Condition {
	return &Condition{Query: field + " like ?", Value: []any{"%" + value + "%"}}
}

func In(field string, value any) *Condition {
	return &Condition{Query: field + " in ?", Value: []any{value}}
}

func NotIn(field string, value any) *Condition {
	return &Condition{Query: field + " not in ?", Value: []any{value}}
}
