package op

import "gorm.io/gorm"

type Clause func(db *gorm.DB) *gorm.DB

func DefaultClause() Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func Limit(limit int) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func Offset(offset int) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func Page(page, pageSize int) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}
}

func Order(clauses ...string) Clause {
	return func(db *gorm.DB) *gorm.DB {
		for _, clause := range clauses {
			db = db.Order(clause)
		}
		return db
	}
}

func Filter(conditions ...*Condition) Clause {
	return func(db *gorm.DB) *gorm.DB {
		for _, condition := range conditions {
			db = db.Where(condition.Query, condition.Value...)
		}
		return db
	}
}

func FilterMap(conditions map[string]any) Clause {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(conditions)
	}
}
