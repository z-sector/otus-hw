package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

const (
	defaultMaxPoolSize  = 2
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Postgres struct {
	Builder squirrel.StatementBuilderType
	Pool    PgPoolI
}

func NewPg(pool PgPoolI) *Postgres {
	return &Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		Pool:    pool,
	}
}

func NewPgByURL(log logger.AppLog, url string) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(defaultMaxPoolSize)
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	var pool PgPoolI
	connAttempts := defaultConnAttempts
	for connAttempts > 0 {
		pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			if err = pool.Ping(context.Background()); err == nil {
				break
			}
		}

		log.Info(fmt.Sprintf("Postgres is trying to connect, attempts left: %d", connAttempts))

		time.Sleep(defaultConnTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return NewPg(pool), nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
