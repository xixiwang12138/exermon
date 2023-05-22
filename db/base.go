package db

import (
	"context"
	"github.com/xixiwang12138/exermon/db/op"
	"gorm.io/gorm"
)

type Transaction *gorm.DB

type BaseDao[T any] struct {
	model *T
	g     *gorm.DB
	ctx   context.Context
}

func NewBaseDao[T any](g *gorm.DB) *BaseDao[T] {
	return &BaseDao[T]{g: g, model: new(T)}
}

func (repo *BaseDao[T]) copy() *gorm.DB {
	return repo.g.Model(repo.model).WithContext(repo.ctx)
}

func (repo *BaseDao[T]) wrap(filter ...*op.Condition) *gorm.DB {
	temp := repo.copy()
	clause := op.Filter(filter...)
	temp = clause(temp)
	return temp
}

// region 事务

func (repo *BaseDao[T]) Begin() *BaseDao[T] {
	return &BaseDao[T]{g: repo.g.Begin()}
}

func (repo *BaseDao[T]) GetTransaction() Transaction {
	return repo.g
}

func (repo *BaseDao[T]) Commit() error {
	return repo.g.Commit().Error
}

func (repo *BaseDao[T]) Rollback() error {
	return repo.g.Rollback().Error
}

// Extend 由事务对象继承出另一个类型的repo
func Extend[R any](ctx context.Context, tx Transaction) *BaseDao[R] {
	return &BaseDao[R]{g: tx, model: new(R), ctx: ctx}
}

// endregion

func (repo *BaseDao[T]) Instance(ctx context.Context) *BaseDao[T] {
	return &BaseDao[T]{
		g:   repo.g,
		ctx: ctx,
	}
}

// region CRUD api

func (repo *BaseDao[T]) Insert(t *T) (err error) {
	temp := repo.copy()
	err = temp.Create(t).Error
	return
}

func (repo *BaseDao[T]) Save(t *T) (err error) {
	temp := repo.copy()
	err = temp.Save(t).Error
	return
}

func (repo *BaseDao[T]) List(clauses ...op.Clause) (list []*T, err error) {
	temp := repo.copy()
	for _, clause := range clauses {
		temp = clause(temp)
	}
	err = temp.Find(list).Error
	return
}

func (repo *BaseDao[T]) Get(filter ...*op.Condition) (r *T, err error) {
	r = new(T)
	temp := repo.wrap(filter...)
	err = temp.First(&r).Error
	return
}

func (repo *BaseDao[T]) Count(filter ...*op.Condition) (count int64, err error) {
	temp := repo.wrap(filter...)
	err = temp.Count(&count).Error
	return
}

func (repo *BaseDao[T]) Delete(filter ...*op.Condition) (err error) {
	temp := repo.wrap(filter...)
	err = temp.Delete(repo.model).Error
	return
}

func (repo *BaseDao[T]) Update(v map[string]any, filter ...*op.Condition) (err error) {
	temp := repo.wrap(filter...)
	err = temp.Updates(v).Error
	return
}

// endregion
