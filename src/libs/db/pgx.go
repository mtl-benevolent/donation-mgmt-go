package db

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/system/logging"
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

	pool, err = pgxpool.NewWithConfig(gs.AppContext(), dbConfig)
	if err != nil {
		l.Error("Unable to connect to the database", slog.Any("error", err))
		panic(fmt.Sprintf("Unable to connect to the database: %v", err))
	}

	l.Info(
		"Database is bootstrapped",
		slog.String("connection_string", connectionString),
		slog.String("host", appConfig.DBHost),
		slog.Any("port", appConfig.DBPort),
		slog.String("database_name", appConfig.AppName),
		// slog.String("app_name", appConfig.AppName),
		slog.String("schema", appConfig.DBSchema),
	)

	rc.RegisterPollComponent(componentName, func() bool {
		l.Debug("Pinging the database")

		pingCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := pool.Ping(pingCtx)

		if err != nil {
			l.Error("Unable to ping the database", slog.Any("error", err))
		}

		return err == nil
	}, pollInterval)

	err = gs.RegisterComponentWithFn(componentName, func() error {
		l.Info("Closing the database connection pool")
		pool.Close()
		return nil
	})
	if err != nil {
		l.Error("Unable to register the database pool with the graceful shutdown", slog.Any("error", err))
		panic("error bootstrapping the database")
	}
}

func BootstrapSingleConnection(appConfig *config.AppConfiguration) (*pgx.Conn, error) {
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

	dbConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		l.Error("Unable to parse database connection string", slog.Any("error", err))
		return nil, err
	}

	dbConfig.Tracer = &queryTracer{logger: l}

	dbConn, err := pgx.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		l.Error("Unable to connect to the database", slog.Any("error", err))
		return nil, err
	}

	return dbConn, nil
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
			sql.WriteString(strings.ReplaceAll(line, "\t", " "))

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
	qt.logger.With(logging.ContextLogData(ctx)...).
		Info("SQL Query started", slog.String("sql", strings.TrimSpace(sql.String())), slog.String("txStatus", txStatus), slog.Any("args", data.Args))
	return ctx
}

func (qt *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// TODO: Add contextual data to the logger
	if data.Err != nil {
		qt.logger.With(logging.ContextLogData(ctx)...).
			Error("SQL Query failed", slog.Any("error", data.Err))
	} else {
		qt.logger.With(logging.ContextLogData(ctx)...).
			Info(
				"SQL Query succeeded",
				slog.Int64("rows_affected", data.CommandTag.RowsAffected()),
			)
	}
}
