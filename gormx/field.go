package gormx

import (
	"gorm.io/gorm"
)

type (
	IntField struct {
		field[int]
	}
	Int8Field struct {
		field[int8]
	}
	Int16Field struct {
		field[int16]
	}
	Int32Field struct {
		field[int32]
	}
	Int64Field struct {
		field[int64]
	}
	UintField struct {
		field[uint]
	}
	Uint8Field struct {
		field[uint8]
	}
	Uint16Field struct {
		field[uint16]
	}
	Uint32Field struct {
		field[uint32]
	}
	Uint64Field struct {
		field[uint64]
	}
	Float32Field struct {
		field[float32]
	}
	Float64Field struct {
		field[float64]
	}
	BoolField struct {
		field[bool]
	}
	StringField struct {
		field[string]
	}
)

type field[T comparable] struct {
	Name string
}

func (n field[T]) ColumnName() string {
	return n.Name
}

func (n field[T]) Eq(v T) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(n.Name+" = ?", v)
	}
}

func (n field[T]) NotEq(v T) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(n.Name+" <> ?", v)
	}
}

func (n field[T]) Gt(v T) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(n.Name+" > ?", v)
	}
}

func (n field[T]) Lt(v T) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(n.Name+" < ?", v)
	}
}

func (n field[T]) Gte(v T) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(n.Name+" >= ?", v)
	}
}

func (n field[T]) Lte(v T) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(n.Name+" <= ?", v)
	}
}

func (n field[T]) In(v []T) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(n.Name+" IN (?)", v)
	}
}

func (n field[T]) NotIn(v []T) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(n.Name+" NOT IN (?)", v)
	}
}

// string field support like operation

func (s *StringField) Like(v string) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(s.Name+" LIKE ?", v)
	}
}

func (s *StringField) NotLike(v string) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(s.Name+" NOT LIKE ?", v)
	}
}

func (s *StringField) PrefixLike(v string) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(s.Name+" LIKE ?", v+"%")
	}
}
