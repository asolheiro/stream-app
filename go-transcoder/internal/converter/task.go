package converter

import (
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
	rabbitmqClient *rabbitmq.RabbitClient
	db             *sql.DB
}

func NewVideoConverter(db *sql.DB, rabbitmqClient *rabbitmq.RabbitClient) *VideoConverter {
	return &VideoConverter{
		db:             db,
		rabbitmqClient: rabbitmqClient,
	}
}

type VideoTask struct {
	VideoId   int    `json:"vide_id"`
	VideoPath string `json:"path"`
}

func (vc *VideoConverter) TaskHandler(delivery amqp.Delivery) {
	var task VideoTask

	err := json.Unmarshal(delivery.Body, &task)
	if err != nil {
		vc.logError(task, "failed to unmarshal task", err)
		return
	}

	if IsProcessed(vc.db, task.VideoId) {
		slog.Warn("Video already processed", slog.Int("video_id", task.VideoId))
		delivery.Ack(false)
		return
	}

	err = vc.processVideo(&task)
	if err != nil {
		vc.logError(task, "failed to process video", err)
		return
	}

	err = MarkProcessed(vc.db, task.VideoId)
	if err != nil {
		vc.logError(task, "failed to mark processed", err)
		return
	}

	delivery.Ack(false)
	slog.Info("Video processed", slog.Int("video_id", task.VideoId))
}

func (vc *VideoConverter) processVideo(task *VideoTask) error {
	mergedFile := filepath.Join(task.VideoPath, "merged.mp4")
	mpegDashPath := filepath.Join(task.VideoPath, "mpeg-dash")

	slog.Info("Merging chunks", slog.String("path", task.VideoPath))
	err := vc.mergeChunks(task.VideoPath, mergedFile)
	if err != nil {
		vc.logError(*task, "failed to merge chunks", err)
		return err
	}

	err = os.MkdirAll(mpegDashPath, os.ModePerm)
	if err != nil {
		vc.logError(*task, "failed to mkdir", err)
		return err
	}

	slog.Info("Converting video to mpeg-dash", slog.String("path", task.VideoPath))
	ffmpegCmd := exec.Command(
		"ffmpeg", "-i", mergedFile,
		"-f", "dash", filepath.Join(mpegDashPath, "output.mpd"),
	)

	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		vc.logError(*task, "failed to convert video to mpeg-dash, output :" + string(output), err)
		return err
	}

	slog.Info("Video converted to mpeg-dash", slog.String("path", mpegDashPath))
	slog.Info("Removing mp4 file", slog.String("path", mergedFile))
	err = os.Remove(mergedFile)
	if err != nil {
		vc.logError(*task, "failed to remove mp4 file", err)
		return err
	}
	return nil
}

func (vc *VideoConverter) logError(task VideoTask, message string, err error) {
	errorData := map[string]any{
		"video_id": task.VideoId,
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

func (vc *VideoConverter) extractNumber(fileName string) int {
	re := regexp.MustCompile(`\d+`)
	numStr := re.FindString(filepath.Base(fileName))
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return -1
	}

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
