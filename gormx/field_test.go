package gormx

import "testing"

var TestTable = &testTable{}

type testTable struct {
	ID     Int64Field
	Name   StringField
	Normal BoolField
}

func TestField(t *testing.T) {
}
