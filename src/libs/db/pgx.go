package db

import (
	"donation-mgmt/src/config"
	"fmt"

	"github.com/gretro/go-lifecycle"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Bootstrap(gs *lifecycle.GracefulShutdown, appConfig *config.AppConfiguration) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?application_name=donation_mgmt&search_path=%s",
		appConfig.DBUser,
		appConfig.DBPassword,
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBName,
		appConfig.DBSchema,
	)

	dbpool, err := pgxpool.New(gs.AppContext(), connectionString)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to the database: %v", err))
	}

	pool = dbpool

	gs.RegisterComponentWithFn("Postgres", func() error {
		dbpool.Close()
		return nil
	})
}

func DBPool() *pgxpool.Pool {
	if pool == nil {
		panic("Database not initialized")
	}

	return pool
}
