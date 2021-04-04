package migration

import (
	"fmt"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/gobuffalo/pop/v5"
)

// MigrateUp does migration
func MigrateUp(config *config.Config) error {
	mig, conn := GetMigrator(config)
	defer conn.Close()

	err := mig.Up()
	return err
}

func GetMigrator(config *config.Config) (pop.FileMigrator, *pop.Connection) {
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

func getConn(config *config.Config) *pop.Connection {
	name := config.DB.Name
	user := config.DB.User
	pass := config.DB.Password
	host := config.DB.Host
	port := config.DB.Port

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
