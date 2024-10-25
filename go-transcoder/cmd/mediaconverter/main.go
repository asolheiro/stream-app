package main

import (
	"context"
	"gotranscoder/internal/converter"
	"gotranscoder/internal/database"
	"gotranscoder/internal/utils"
	"gotranscoder/pkg/log"
	"gotranscoder/pkg/rabbitmq"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/streadway/amqp"
)

func main() {
	//VarEnvs
	// Consumer
	conversionExchange := utils.GetEnvOrDefault("CONVERSION_EXCHANGE", "conversionExchange")
	conversionQueue := utils.GetEnvOrDefault("CONVERSION_QUEUE", "videoConversion_queue")
	conversionKey := utils.GetEnvOrDefault("CONVERTION_KEY", "videoConversion")

	//Producer
	confirmationQueue := utils.GetEnvOrDefault("CONFIRMATION_QUEUE", "conversionConfirmation_queue")
	confirmationKey := utils.GetEnvOrDefault("CONFIRMATION_KEY", "videoConfirmation	")
	rootPath := utils.GetEnvOrDefault("VIDEO_ROOT_PATH", "media/uploads")
	//

	isDebug := utils.GetEnvOrDefault("DEBUG", "false") == "true"
	logger := log.NewLogger(isDebug)
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChann := make(chan os.Signal, 1)
	signal.Notify(signalChann, syscall.SIGINT, syscall.SIGTERM)

	db, err := database.ConnectPostgres()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitmqURL := utils.GetEnvOrDefault("RABBITMQ_URL", "ampq://guest:guest@localhost:5672/")
	rabbitmqClient, err := rabbitmq.NewRabbitClient(ctx, rabbitmqURL)
	if err != nil {
		panic(err)
	}
	defer rabbitmqClient.Close()

	videoConverter := converter.NewVideoConverter(db, rabbitmqClient, rootPath)

	messagesChannel, err := rabbitmqClient.ConsumeMessages(conversionExchange, conversionKey, conversionQueue)
	if err != nil {
		slog.Error("failed to consume messages", slog.String("error", err.Error()))
		return
	}

	var wg sync.WaitGroup
	go func() {
		for delivery := range messagesChannel {
			wg.Add(1)	
			go func (delivery amqp.Delivery)  {
				defer wg.Done()
				videoConverter.HandleMessage(ctx, delivery, conversionExchange, confirmationKey, confirmationQueue)
			}(delivery)
		}
	}()
	
	slog.Info("Waiting for messages from RabbitMQ")
	<-signalChann
	slog.Info("Shutdown signal received, finalizing processing...")
	
	cancel()
	wg.Wait()
	slog.Info("Processubg completed, exiting...")
}
