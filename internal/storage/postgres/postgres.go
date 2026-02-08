package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"taskmgr/internal/config"
	"taskmgr/internal/storage"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DBConfig) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx config parse: %w", err)
	}

	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pgx pool create: %w", err)
	}

	pg := &Postgres{pool: pool}

	// Проверяем подключение
	if err := pg.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *Postgres) Close() error {
	p.pool.Close()
	return nil
}

// Убедимся, что реализует интерфейс Storage
var _ storage.Storage = (*Postgres)(nil)
