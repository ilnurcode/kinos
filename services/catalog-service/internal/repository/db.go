// Package repository предоставляет типы для управления подключением к базе данных.
// Включает интерфейс DB и реализацию PostgresDB.
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB interface {
	ConnDB(dsn string) error
	CloseDB() error
}

type PostgresDB struct {
	pool *pgxpool.Pool
}

func (p *PostgresDB) ConnDB(dsn string) error {
	var err error
	p.pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

func (p *PostgresDB) CloseDB() error {
	if p.pool != nil {
		p.pool.Close()
	}
	return nil
}

func (p *PostgresDB) GetPool() *pgxpool.Pool {
	return p.pool
}
