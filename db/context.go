package db

import (
	"context"
	"database/sql"

	"github.com/xixiwang12138/exermon/flow"
)

var (
	DefaultTxLevel = sql.LevelRepeatableRead
)

var txCtxKey = struct{}{}

func MustTransaction(ctx context.Context, task func(context.Context)) {
	gormDB := Component.Gorm()

	tx := gormDB.Begin(&sql.TxOptions{
		Isolation: DefaultTxLevel,
	})

	ctx = context.WithValue(ctx, txCtxKey, tx)

	flow.Try(func() {
		task(ctx)
		tx.Commit()
	}, func(err interface{}) {
		tx.Rollback()
		panic(err)
	})
}
