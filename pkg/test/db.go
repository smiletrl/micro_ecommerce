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
	contextd "github.com/smiletrl/micro_ecommerce/pkg/context"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	connStr := cfg.DBConnString
	if cfg.DBConnString == "" {
		connStr = fmt.Sprintf("user=%s sslmode=%s host=%s password=%s port=%s dbname=%s", cfg.DB.User, cfg.DB.SSLMode, cfg.DB.Host, cfg.DB.Password, cfg.DB.Port, cfg.DB.Name)
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

// Context returns mock context
func Context() *contextd.Context {
	e := echo.New()
	req := &http.Request{}
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	cc := contextd.NewMock(c)
	return cc
}
