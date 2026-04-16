// Package repository предоставляет транзакционный менеджер для inventory-service.
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManagerInterface interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type txKey struct{}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

func (tm *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if existingTx := ctx.Value(txKey{}); existingTx != nil {
		return fn(ctx)
	}

	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = commitErr
			}
		}
	}()

	txCtx := context.WithValue(ctx, txKey{}, tx)
	err = fn(txCtx)
	return err
}
