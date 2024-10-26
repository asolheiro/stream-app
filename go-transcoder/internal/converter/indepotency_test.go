//go:build testcontainers
// +build testTag

package converter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"gotranscoder/internal/converter"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIsProcessed(t *testing.T) {
	ctx := context.Background()

	postgresContainer, db, err := SetupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("failed to setup PostgreSQ container: %v", err)
	}
	defer postgresContainer.Terminate(ctx)
	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO processed_videos (video_id, status, processed_at) VALUES ($1, $2, $3)",
		1, "success", time.Now(),
	)
	assert.NoError(t, err, "failed to insert ")

	isProcessed := converter.IsProcessed(db, 1)
	assert.True(t, isProcessed, "isProcessed should return true")

	isProcessed = converter.IsProcessed(db, 999)
	assert.False(t, isProcessed, "isProcessed should return false")
}

func TestMarkProcessed(t *testing.T) {
	ctx := context.Background()

	postgresC, db, err := SetupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("failed to setup PostgresSQL container: %v", err)
	}
	defer postgresC.Terminate(ctx)
	defer db.Close()

	err = converter.MarkProcessed(db, 2)
	assert.NoError(t, err, "mark processed should not return an error")

	var status string
	err = db.QueryRow(
		"SELECT status FROM processed_videos WHERE video_id = $1",
		2,
	).Scan(&status)
	assert.NoError(t, err, "query should return a value")
	assert.Equal(t, "success", status, "db status should be 'success'")
}

func TestRegisterError(t *testing.T) {
	ctx := context.Background()

	postgresC, db, err := SetupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("failed to setup PostgreSQL container: %v" )
	}
	defer postgresC.Terminate(ctx)
	defer db.Close()

	errorData := map[string]any{
		"video_id": 1,
		"phases": []string{
			"Phase1", "Phase2",
		},
		"error_msg": "Test_error",
	}

	converter.RegisterError(db, errorData, fmt.Errorf("Test error"))

	var errorDetails []byte
	err = db.QueryRow(
		"SELECT error_details FROM process_errors_log WHERE id = 1",
	).Scan(&errorDetails)
	assert.NoError(t, err, "query should return a value")

	var loggedError map[string]interface{}
	err = json.Unmarshal(errorDetails, &loggedError)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), loggedError["video_id"].(float64))
	assert.Equal(t, "Test error", loggedError["error_msg"])
}
