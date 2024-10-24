package database

import (
	"database/sql"
	"fmt"
	"gotranscoder/internal/utils"
	"log/slog"

	_ "github.com/lib/pq"
)

func ConnectPostgres() (*sql.DB, error) {
	user := utils.GetEnvOrDefault("POSTGRES_USER", "root")
	password := utils.GetEnvOrDefault("POSTGRES_PASSWORD", "password")
	dbname := utils.GetEnvOrDefault("POSTGRES_DB", "converter_database")
	host := utils.GetEnvOrDefault("POSTGRES_HOST", "postgres")
	sslMode := utils.GetEnvOrDefault("POSTGRES_SSL_MODE", "disable")

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
