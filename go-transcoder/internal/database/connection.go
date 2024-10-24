package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

func ConnectPostgres() (*sql.DB, error) {
	user := getEnvOrDefault("POSTGRES_USER", "root")
	password := getEnvOrDefault("POSTGRES_PASSWORD", "password")
	dbname := getEnvOrDefault("POSTGRES_DB", "converter_database")
	host := getEnvOrDefault("POSTGRES_HOST", "postgres")
	sslMode := getEnvOrDefault("POSTGRES_SSL_MODE", "disable")

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=%s",
		user, password, dbname, host, sslMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("error connecting to database", slog.String("connStr", connStr))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Error("error pinging database", slog.String("connStr", connStr))
		return nil, err
	}
	slog.Info("connected to database successfully")
	return db, nil
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
