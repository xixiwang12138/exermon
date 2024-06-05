package gormx

import "gorm.io/gorm"

type Clause func(db *gorm.DB) *gorm.DB

type Column interface {
	ColumnName() string
}
