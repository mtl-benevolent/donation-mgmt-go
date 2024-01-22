package main

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/permissions"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

var roCapabilities = []string{
	permissions.EntityOrganization.Capability(permissions.ActionRead),
	permissions.EntityDonation.Capability(permissions.ActionRead),
}

var userCapabilities = append(
	roCapabilities,
	permissions.EntityDonation.Capability(permissions.ActionCreate),
	permissions.EntityDonation.Capability(permissions.ActionUpdate),
)

var managerCapabilities = append(
	userCapabilities,
	permissions.EntityOrganization.Capability(permissions.ActionUpdate),
	permissions.EntityDonation.Capability(permissions.ActionDelete),
	permissions.EntityRoles.Capability(permissions.ActionRead),
	permissions.EntityRoles.Capability(permissions.ActionUpdate),
)

var adminCapabilities = append(
	managerCapabilities,
	permissions.EntityOrganization.Capability(permissions.ActionCreate),
)

func main() {
	appConfig := config.Bootstrap()
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.ForceSetLogger(l)

	defer func() {
		err := recover()
		if err != nil {
			l.Error("Failed at seeding DB", slog.Any("error", err))

			os.Exit(1)
		}
	}()

	pgConn, err := db.BootstrapSingleConnection(appConfig)
	if err != nil {
		panic("Unable to connect to the database: " + err.Error())
	}
	defer pgConn.Close(context.Background())

	tx, err := pgConn.Begin(context.Background())
	if err != nil {
		panic("Unable to start transaction: " + err.Error())
	}

	err = seedRoles(l, tx)
	if err != nil {
		rollback(tx)
		os.Exit(1)
		return
	}

	err = seedGlobalRoles(l, tx)
	if err != nil {
		rollback(tx)
		os.Exit(1)
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		panic("Unable to commit transaction: " + err.Error())
	}

	l.Info("Successfully seeded DB")
}

func seedRoles(logger *slog.Logger, tx pgx.Tx) error {
	logger.Info("Seeding roles...")

	_, err := tx.Exec(context.Background(),
		`INSERT INTO roles (name, capabilities) 
		VALUES 
			($1, $2),
			($3, $4),
		  ($5, $6),
		 	($7, $8)
		ON CONFLICT (name) DO UPDATE SET capabilities = EXCLUDED.capabilities`,
		"read-only", roCapabilities,
		"user", userCapabilities,
		"manager", managerCapabilities,
		"admin", adminCapabilities,
	)

	if err != nil {
		logger.Error("Failed to seed roles", slog.Any("error", err))
		return err
	}

	logger.Info("Seeded roles successfully")
	return nil
}

func seedGlobalRoles(logger *slog.Logger, tx pgx.Tx) error {
	logger.Info("Seeding global user roles...")

	_, err := tx.Exec(context.Background(),
		`insert into global_user_roles(subject, role_id)
		select $1 as subject, id as role_id from roles r where r."name" = $2
		on conflict do nothing`,
		"root", "admin",
	)

	if err != nil {
		logger.Error("Failed to seed global user roles", slog.Any("error", err))
		return err
	}

	logger.Info("Seeded global user roles successfully")
	return nil
}

func rollback(tx pgx.Tx) {
	err := tx.Rollback(context.Background())
	if err != nil {
		panic("Unable to rollback transaction: " + err.Error())
	}
}
