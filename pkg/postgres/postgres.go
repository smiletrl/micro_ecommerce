package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/migration"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
)

type Provider interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)

	// used at defer after provider initialization.
	Close()
}

type provider struct {
	db      *pgxpool.Pool
	tracing tracing.Provider
}

// NewProvider returns a new postgresql db
func NewProvider(cfg config.Config, tracing tracing.Provider) (Provider, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgresql.User, cfg.Postgresql.Password, cfg.Postgresql.Host, cfg.Postgresql.Port, cfg.Postgresql.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dbpool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := migration.MigrateUp(cfg); err != nil {
		dbpool.Close()

		return nil, err
	}

	return provider{dbpool, tracing}, nil
}

func (p provider) RawDB() *pgxpool.Pool {
	return p.db
}

func (p provider) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	span, ctx := p.tracing.StartSpan(ctx, sql)
	defer p.tracing.FinishSpan(span)

	return p.db.Query(ctx, sql, args...)
}

func (p provider) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	span, ctx := p.tracing.StartSpan(ctx, sql)
	defer p.tracing.FinishSpan(span)

	return p.db.QueryRow(ctx, sql, args...)
}

func (p provider) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	span, ctx := p.tracing.StartSpan(ctx, sql)
	defer p.tracing.FinishSpan(span)

	return p.db.Exec(ctx, sql, args...)
}

func (p provider) Close() {
	p.db.Close()
}
