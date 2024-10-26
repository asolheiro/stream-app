package log_test

import (
	"bytes"
	"gotranscoder/pkg/log"
	"io"
	"os"
	"testing"
)

func TestLoggerOutput(t *testing.T) {
	// Create Pipe to capture the output
	r, w, _ := os.Pipe()
	stdOut := os.Stdout
	os.Stdout = w

	// Create a logger and log a message
	log := log.NewLogger(false)
	log.Error("Testing logger output")

	// Close the writer and restore os.Stdout
	w.Close()
	os.Stdout = stdOut

	var bufffer bytes.Buffer
	io.Copy(&bufffer, r)

	// Verify that log message was written to os.Stdout
	if !containsLogLevel(bufffer.String(), "ERROR") {
		t.Errorf("Expected log level ERROR, but got: %s", bufffer.String())
	}
}

func containsLogLevel(output, level string) bool {
	return bytes.Contains([]byte(output), []byte(level))
}
