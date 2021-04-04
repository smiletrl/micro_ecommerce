package dbcontext

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/migration"
)

// DB adds some wrappers around standard sqlx functionality
type DB interface {
	Query(c echo.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(c echo.Context, sql string, args ...interface{}) pgx.Row
	Exec(c echo.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type dbcontext struct {
	DB *pgxpool.Pool
}

// InitDB is to inti db
func InitDB(cfg config.Config) (DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	dbpool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err := migration.MigrateUp(cfg); err != nil {
		panic(err)
	}

	return NewDBContext(dbpool), nil
}

// NewDBContext returns a new eps db wrapper around an sqlx.DB
func NewDBContext(db *pgxpool.Pool) DB {
	return &dbcontext{db}
}

func (db *dbcontext) RawDB() *pgxpool.Pool {
	return db.DB
}

func (db *dbcontext) Query(c echo.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	// @todo add tracing later to monitor the performance
	return db.DB.Query(c.Request().Context(), sql, args...)
}

func (db *dbcontext) QueryRow(c echo.Context, sql string, args ...interface{}) pgx.Row {
	// @todo add tracing later to monitor the performance
	return db.DB.QueryRow(c.Request().Context(), sql, args...)
}

func (db *dbcontext) Exec(c echo.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	// @todo add tracing later to monitor the performance
	return db.DB.Exec(c.Request().Context(), sql, args...)
}
