package migration

import (
	"fmt"
	"github.com/gobuffalo/pop/v5"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
)

// MigrateUp does migration
func MigrateUp(config config.Config) error {
	mig, conn := GetMigrator(config)
	defer conn.Close()

	err := mig.Up()
	return err
}

func GetMigrator(config config.Config) (pop.FileMigrator, *pop.Connection) {
	// default path defined in ./bin/docker-entrypoint.sh
	migrationPath := "/app/migrations"
	if config.MigrationPath != "" {
		migrationPath = config.MigrationPath
	}
	conn := getConn(config)
	mig, err := pop.NewFileMigrator(migrationPath, conn)
	if err != nil {
		panic(err)
	}
	return mig, conn
}

func getConn(config config.Config) *pop.Connection {
	name := config.Postgresql.Name
	user := config.Postgresql.User
	pass := config.Postgresql.Password
	host := config.Postgresql.Host
	port := config.Postgresql.Port

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, name)

	cd := &pop.ConnectionDetails{
		URL: url,
	}
	con, err := pop.NewConnection(cd)
	if err != nil {
		panic(err)
	}
	return con
}
