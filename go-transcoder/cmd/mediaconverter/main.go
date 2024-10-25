package main

import (
	"gotranscoder/internal/converter"
	"gotranscoder/internal/database"
	"gotranscoder/internal/rabbitmq"
	"gotranscoder/internal/utils"
	"log/slog"

	"github.com/streadway/amqp"
)

func main() {
	//VarEnvs
	// Consumer
	exchangeName := utils.GetEnvOrDefault("CONVERSION_EXCHANGE", "conversionExchange")
	queueName := utils.GetEnvOrDefault("CONVERSION_QUEUE", "videoConversion_queue")
	routingKey := utils.GetEnvOrDefault("CONVERTION_KEY", "videoConversion")

	//Producer
	confirmationExc := utils.GetEnvOrDefault("CONFIRMATION_EXCHANGE", "confirmationExchange")
	confirmationQueue := utils.GetEnvOrDefault("CONFIRMATION_QUEUE", "conversionConfirmation_queue")
	confirmationKey := utils.GetEnvOrDefault("CONFIRMATION_KEY", "videoConfirmation	")

	
	
	db, err := database.ConnectPostgres()
	if err != nil {
		slog.Error("failed to connect to PostgreSQL", slog.String("error", err.Error()))
		panic(err)
	}

	rabbitmqURL := utils.GetEnvOrDefault("RABBITMQ_URL", "ampq://guest:guest@localhost:5672/")
	rabbitmqClient, err := rabbitmq.NewRabbitClient(rabbitmqURL)
	if err != nil {
		panic(err)
	}
	defer rabbitmqClient.Close()

	vc := converter.NewVideoConverter(db, rabbitmqClient)

	messagesChannel, err := rabbitmqClient.ConsumeMessages(exchangeName, routingKey, queueName)
	if err != nil {
		slog.Error("failed to consume messages", slog.String("error", err.Error()))
	}

	for messageDelivered := range messagesChannel {
		go func(delivery amqp.Delivery) {
			vc.HandleMessage(
				context,
				delivery,
				confirmationExc,
				confirmationKey,
				confirmationQueue,
			)
		}(messageDelivered)

	}


}
