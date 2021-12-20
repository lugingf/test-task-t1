package dat

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var dbMap *pgxpool.Pool

// InitDB will initialise a connection to the database and return it
func InitDB(ctx context.Context, logger *zap.Logger) (*pgxpool.Pool, error) {

	dbHost := envVar("DB_HOST", "localhost")
	dbUser := envVar("DB_USER", "talon")
	dbPassword := envVar("DB_PASSWORD", "talon.one.8080")
	dbName := envVar("DB_NAME", "talon")
	dbPort := envVar("DB_PORT", "5432")
	dbSSLMode := envVar("DB_SSL", "disable")

	if dbHost == "" {
		return nil, errors.New("the database hostname may not be empty. Please provide a hostname by setting the environment variable DB_HOST")
	} else if dbUser == "" {
		return nil, errors.New("the database username may not be empty. Please provide a username by setting the environment variable DB_USER")
	} else if dbPassword == "" {
		return nil, errors.New("the database password may not be empty. Please provide a password by setting the environment variable DB_PASSWORD")
	} else if dbName == "" {
		return nil, errors.New("the database name may not be empty. Please provide the name of the database by setting the environment variable DB_NAME")
	} else if dbPort == "" {
		return nil, errors.New("the port for the database connection may not be empty. Please provide the port by setting the environment variable DB_PORT")
	} else if dbSSLMode == "" {
		logger.Warn("No SSL mode for the database connection was provided. Consider setting the SSL mode by setting the environement variable DB_SSL")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Building pgx connConfig")
	}

	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		logger.Info("Connected to Database", zap.String("host", dbHost), zap.String("name", dbName), zap.String("user", dbUser))
		return nil
	}

	connConfig.MaxConnIdleTime = time.Minute * 30

	hostname, err := os.Hostname()
	if err == nil && hostname != "" {
		connConfig.ConnConfig.RuntimeParams["application_name"] = hostname
	}
	connConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	connConfig.MaxConns = 100

	db, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Connection pool")
	}

	dbMap = db
	return db, nil
}

func envVar(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return defaultValue
}

// ExitDb closes the open connections to the database
func ExitDb(logger *zap.Logger) {
	dbMap.Close()
	logger.Info("closed connection to DB")
}
