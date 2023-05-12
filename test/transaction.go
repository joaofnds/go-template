package test

import (
	"context"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var TransactionModule = fx.Module(
	"transaction",
	fx.Decorate(BeginTransaction),
	fx.Invoke(RollbackTransaction),
)

func BeginTransaction(db *gorm.DB) *gorm.DB {
	return db.Begin()
}

func RollbackTransaction(lifecycle fx.Lifecycle, db *gorm.DB) {
	lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return db.Rollback().Error
		},
	})
}
