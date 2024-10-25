package converter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gotranscoder/internal/rabbitmq"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type VideoConverter struct {
	db             *sql.DB
	rabbitmqClient *rabbitmq.RabbitClient
	rootPath       string
}

func NewVideoConverter(db *sql.DB, rabbitmqClient *rabbitmq.RabbitClient, rootPath string) *VideoConverter {
	return &VideoConverter{
		db:             db,
		rabbitmqClient: rabbitmqClient,
		rootPath: rootPath,
	}
}

type VideoTask struct {
	VideoID   int    `json:"vide_id"`
	VideoPath string `json:"path"`
}

func (vc *VideoConverter) HandleMessage(ctx context.Context, delivery amqp.Delivery, confirmationExc, confirmationKey, confirmationQueue string) {
	var task VideoTask

	if err := json.Unmarshal(delivery.Body, &task); err != nil {
		vc.logError(task, "failed to deserialize message", err)
		delivery.Ack(false)
		return
	}

	if IsProcessed(vc.db, task.VideoID) {
		slog.Warn("Video already processed", slog.Int("video_id", task.VideoID))
		delivery.Ack(false)
		return
	}

	if err := vc.processVideo(&task); err != nil {
		vc.logError(task, "failed to process video", err)
		delivery.Ack(false)
		return
	}
	slog.Info("Video conversion processed", slog.Int("video_id", task.VideoID))

	if err := MarkProcessed(vc.db, task.VideoID); err != nil {
		vc.logError(task, "failed to mark processed", err)
		return
	}
	delivery.Ack(false)
	slog.Info("Video marked as processed", slog.Int("video_id", task.VideoID))

	confirmationMessage := []byte(
		fmt.Sprintf(`{"video_id": %d, "path": "%s"}`, task.VideoID, task.VideoPath),
	)

	if err := vc.rabbitmqClient.PublishMessage(confirmationExc, confirmationKey, confirmationQueue, confirmationMessage); err != nil {
		vc.logError(task, "failed to publish confirmation message", err)
	}
	slog.Info("Confirmation message published on", slog.String("queue", confirmationQueue), slog.Int("video_id", task.VideoID))
}

func (vc *VideoConverter) processVideo(task *VideoTask) error {
	chunkPath := filepath.Join(vc.rootPath, fmt.Sprintf("%d", task.VideoID))
	mergedFile := filepath.Join(task.VideoPath, "merged.mp4")
	mpegDashPath := filepath.Join(task.VideoPath, "mpeg-dash")

	slog.Info("Merging chunks", slog.String("path", chunkPath))
	if err := vc.mergeChunks(chunkPath, mergedFile); err != nil {
		return fmt.Errorf("failed to merge chunks")
	}

	if err := os.MkdirAll(mpegDashPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	slog.Info("Converting video to mpeg-dash", slog.String("path", task.VideoPath))
	ffmpegCmd := exec.Command(
		"ffmpeg", "-i", mergedFile,
		"-f", "dash", filepath.Join(mpegDashPath, "output.mpd"),
	)

	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to convert to MPEG-DASH: %v, output: %s", err, string(output))
	}
	slog.Info("Video converted to MPEG-DASH", slog.String("path", mpegDashPath))

	if err := os.Remove(mergedFile); err != nil {
		slog.Warn("failed to remove merged file", slog.String("file", mergedFile), slog.String("error", err.Error()))
		return err
	}
	slog.Info("MP4 merged file removed", slog.String("file", mergedFile))
	return nil
}


func (vc *VideoConverter) extractNumber(fileName string) int {
	re := regexp.MustCompile(`\d+`)
	numStr := re.FindString(filepath.Base(fileName))
	num, _ := strconv.Atoi(numStr)
	return num
}

func (vc *VideoConverter) mergeChunks(inputDir, outputFile string) error {
	chunks, err := filepath.Glob(filepath.Join(inputDir, "*.chunk"))
	if err != nil {
		return fmt.Errorf("failed to find chunks: %v", err)
	}
	
	sort.Slice(chunks, func(i, j int) bool {
		return vc.extractNumber(chunks[i]) < vc.extractNumber(chunks[j])
	})
	
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()
	
	for i, chunk := range chunks {
		input, err := os.Open(chunk)
		if err != nil {
			return fmt.Errorf("failed to open chunk %d: %v", i+1, chunk)
		}
		
		_, err = output.ReadFrom(input)
		if err != nil {
			return fmt.Errorf("failed to write chunk %s to merged file: %v", chunk, err)
		}
		input.Close()
	}
	return nil
}

func (vc *VideoConverter) logError(task VideoTask, message string, err error) {
	errorData := map[string]any{
		"video_id": task.VideoID,
		"error":    message,
		"details":  err.Error(),
		"time":     time.Now(),
	}
	serializedError, _ := json.Marshal(errorData)
	slog.Error("processing err", slog.String("error_details", string(serializedError)))

	regErr := RegisterError(vc.db, errorData, err)
	if regErr != nil {
		slog.Error("failed to register err", slog.String("error_details", string(serializedError)))
	}
}
