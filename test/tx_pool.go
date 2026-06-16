package test

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TxPool wraps a *sql.DB and lets a test open one transaction and roll it back,
// undoing everything the test wrote — without rebuilding the fx app per test.
// The harness opens the transaction in BeforeEach and rolls it back in AfterEach.
//
// A bare db.Exec("BEGIN")/("ROLLBACK") on a pooled *gorm.DB does not isolate:
// database/sql only treats a transaction as real when it is opened through
// BeginTx, and pgx discards a pooled connection returned to it mid-transaction.
// Routing queries through an actual *sql.Tx is what makes the rollback take.
type TxPool struct {
	db *sql.DB
	mu sync.RWMutex
	tx *sql.Tx
}

func (pool *TxPool) Begin() error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.db == nil {
		return errors.New("tx pool: not attached to a database")
	}
	if pool.tx != nil {
		return errors.New("tx pool: transaction already open")
	}

	tx, err := pool.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	pool.tx = tx
	return nil
}

func (pool *TxPool) Rollback() error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.tx == nil {
		return nil
	}

	err := pool.tx.Rollback()
	pool.tx = nil
	return err
}

// active returns the open transaction when there is one, the pool otherwise.
func (pool *TxPool) active() gorm.ConnPool {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	if pool.tx != nil {
		return pool.tx
	}
	return pool.db
}

// gormConnPool is the ConnPool gorm sees. It routes every query to the TxPool's
// active connection. It implements gorm's TxCommitter (with no-op commit and
// rollback) so gorm treats the database as already inside a transaction: code
// that opens its own gorm transaction — e.g. the casbin adapter's
// db.Transaction — then reuses the test transaction instead of starting a new
// one. The per-test rollback on TxPool stays the single source of truth.
type gormConnPool struct{ pool *TxPool }

func (c gormConnPool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return c.pool.active().PrepareContext(ctx, query)
}

func (c gormConnPool) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.pool.active().ExecContext(ctx, query, args...)
}

func (c gormConnPool) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.pool.active().QueryContext(ctx, query, args...)
}

func (c gormConnPool) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return c.pool.active().QueryRowContext(ctx, query, args...)
}

// Commit and Rollback are no-ops: a nested gorm transaction must not end the
// test transaction. TxPool.Rollback (called by the harness) owns that.
func (c gormConnPool) Commit() error   { return nil }
func (c gormConnPool) Rollback() error { return nil }

// GetDBConn lets gorm's db.DB() reach the real pool (health checks, Ping, Close)
// rather than the transaction.
func (c gormConnPool) GetDBConn() (*sql.DB, error) { return c.pool.db, nil }

// TxIsolation rewires the *gorm.DB so every query flows through a shared TxPool,
// and exposes that pool so the harness can open and roll back one transaction
// per test. PrepareStmt is intentionally left off here: a cached *sql.Stmt would
// bind to the transaction it was first prepared on and break after rollback.
var TxIsolation = fx.Options(
	fx.Provide(func() *TxPool { return &TxPool{} }),
	fx.Decorate(func(db *gorm.DB, pool *TxPool) (*gorm.DB, error) {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		pool.db = sqlDB

		return gorm.Open(
			postgres.New(postgres.Config{Conn: gormConnPool{pool: pool}}),
			&gorm.Config{
				Logger:                   db.Logger,
				SkipDefaultTransaction:   true,
				DisableNestedTransaction: true,
			},
		)
	}),
)
