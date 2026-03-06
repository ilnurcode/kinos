// Package repository предоставляет менеджер транзакций для выполнения операций в базе данных.
// Включает интерфейс TransactionManager и реализацию TxManager.
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type txKey struct{}

type TxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) *TxManager {
	return &TxManager{db: db}
}

func (tm *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {

	if existingTx := ctx.Value(txKey{}); existingTx != nil {
		return fn(ctx)
	}
	tx, err := tm.db.Begin(ctx)
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
			if cerr := tx.Commit(ctx); cerr != nil {
				err = cerr
			}
		}
	}()

	txCtx := context.WithValue(ctx, txKey{}, tx)

	err = fn(txCtx)
	return err
}
