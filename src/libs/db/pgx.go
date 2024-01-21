package db

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/system/contextual"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gretro/go-lifecycle"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const componentName = "PostgreSQL"
const pollInterval = 30 * time.Second

var pool *pgxpool.Pool

func Bootstrap(gs *lifecycle.GracefulShutdown, rc *lifecycle.ReadyCheck, appConfig *config.AppConfiguration) {
	l := logger.ForComponent(componentName)

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?application_name=%s&search_path=%s",
		appConfig.DBUser,
		appConfig.DBPassword,
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBName,
		appConfig.AppName,
		appConfig.DBSchema,
	)

	dbConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		l.Error("Unable to parse database connection string", slog.Any("error", err))
		panic("error bootstrapping the database")
	}

	dbConfig.ConnConfig.Tracer = &queryTracer{logger: l}

	dbpool, err := pgxpool.NewWithConfig(gs.AppContext(), dbConfig)
	if err != nil {
		l.Error("Unable to connect to the database", slog.Any("error", err))
		panic(fmt.Sprintf("Unable to connect to the database: %v", err))
	}

	pool = dbpool

	l.Info("Database is bootstrapped")

	rc.RegisterPollComponent(componentName, func() bool {
		l.Debug("Pinging the database")

		pingCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := dbpool.Ping(pingCtx)

		if err != nil {
			l.Error("Unable to ping the database", slog.Any("error", err))
		}

		return err == nil
	}, pollInterval)

	gs.RegisterComponentWithFn(componentName, func() error {
		l.Info("Closing the database connection pool")
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

type queryTracer struct {
	logger *slog.Logger
}

func (qt *queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	lines := strings.Split(data.SQL, "\n")
	sql := strings.Builder{}

	for i, line := range lines {
		if !strings.HasPrefix(line, "--") {
			sql.WriteString(line)

			if i > 0 {
				sql.WriteString(" ")
			}
		}
	}

	txStatus := "unknown"
	switch conn.PgConn().TxStatus() {
	case 'I':
		txStatus = "idle"
	case 'T':
		txStatus = "in_transaction"
	case 'E':
		txStatus = "failed"
	}

	// TODO: Add Contextual data to the logger
	qt.logger.With(contextual.ContextLogData(ctx)...).
		Info("SQL Query started", slog.String("sql", strings.TrimSpace(sql.String())), slog.String("txStatus", txStatus), slog.Any("args", data.Args))
	return ctx
}

func (qt *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// TODO: Add contextual data to the logger
	if data.Err != nil {
		qt.logger.With(contextual.ContextLogData(ctx)...).
			Error("SQL Query failed", slog.Any("error", data.Err))
	} else {
		qt.logger.With(contextual.ContextLogData(ctx)...).
			Info(
				"SQL Query succeeded",
				slog.Int64("rows_affected", data.CommandTag.RowsAffected()),
			)
	}
}
