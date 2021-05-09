package test

import (
	"context"
	"fmt"
	"github.com/gobuffalo/pop/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	configd "github.com/smiletrl/micro_ecommerce/pkg/config"

	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

var (
	dbpool   *pgxpool.Pool
	db       *sqlx.DB
	cfg      configd.Config
	migrator pop.FileMigrator
	m        *migrate.Migrate
)

func init() {
	var err error
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..") + "/.."
	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	stage := os.Getenv("STAGE")
	if stage == "" {
		stage = "local"
	}
	cfg, err = configd.Load(stage)
	if err != nil {
		panic(err)
	}

	// Only use this extra config DBConnString to make it work for github action test.
	// Not sure why, but github action jobs can not parse the following fmt sprintf string.
	connStr := cfg.PostgresqlConnString
	if cfg.PostgresqlConnString == "" {
		connStr = fmt.Sprintf("user=%s sslmode=%s host=%s password=%s port=%s dbname=%s", cfg.Postgresql.User, cfg.Postgresql.SSLMode, cfg.Postgresql.Host, cfg.Postgresql.Password, cfg.Postgresql.Port, cfg.Postgresql.Name)
	}

	dbpool, err = pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	db, err = sqlx.Connect("pgx", connStr)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err = migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		panic(err)
	}
}

// DB returns the db instance
func DB() (*pgxpool.Pool, configd.Config, error) {
	if err := m.Down(); err != nil && err.Error() != "no change" {
		return nil, cfg, err
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		return nil, cfg, err
	}

	content, err := ioutil.ReadFile("testdata/test_insert.sql")
	if err != nil {
		return nil, cfg, err
	}

	if _, err := db.Exec(string(content)); err != nil {
		return nil, cfg, err
	}

	return dbpool, cfg, nil
}
