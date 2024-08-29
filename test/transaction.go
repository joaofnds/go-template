package test

import (
	"context"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Transaction = fx.Decorate(BeginTransaction)

func BeginTransaction(lifecycle fx.Lifecycle, db *gorm.DB) *gorm.DB {
	tx := db.Begin()
	lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error { return tx.Rollback().Error },
	})
	return tx
}
