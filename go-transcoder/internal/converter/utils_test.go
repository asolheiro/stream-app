package converter_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresUser     = "testuser"
	postgresPassword = "testpassword"
	postgresDB       = "testdb"
)

func startRabbitMQContainer(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image: "rabbitmq:3-management",
		ExposedPorts: []string{
			"5672/tcp",
			"15672/tcp",
		},
		WaitingFor: wait.ForLog("Server startup complete"),
	}

	rabbitmqC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	host, err := rabbitmqC.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	port, err := rabbitmqC.MappedPort(ctx, "5672")
	if err != nil {
		return nil, "", err
	}

	rabbitmqURL := fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port.Port())
	fmt.Println("RabbitMQ URL: ", rabbitmqURL)
	return rabbitmqC, rabbitmqURL, nil
}

func setupPostgresContainer(ctx context.Context) (testcontainers.Container, *sql.DB, error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:17.0-alpine3.20",
		ExposedPorts: []string{
			"5432/tcp",
		},
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDB,
		},
		WaitingFor: wait.ForListeningPort("5672/tcp"),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	host, err := postgresC.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, host, port.Port(), postgresDB)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	sqlFilePath := filepath.Join("../../", "db.sql")
	sqlContent, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return nil, nil, err
	}

	_, err = db.Exec(string(sqlContent))
	if err != nil {
		return nil, nil, err
	}

	return postgresC, db, nil
}
