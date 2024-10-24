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
	db, err := database.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	rabbitmqURL := utils.GetEnvOrDefault("RABBITMQ_URL", "ampq://guest:guest@localhost:5672/")
	rabbitmqClient, err := rabbitmq.NewRabbitClient(rabbitmqURL)
	if err != nil {
		panic(err)
	}
	defer rabbitmqClient.Close()

	exchangeName := utils.GetEnvOrDefault("CONVERSION_EXCHANGE", "conversionExchange")
	queueName:= utils.GetEnvOrDefault("CONVERSION_QUEUE", "videoConversion_queue")
	routingKey := utils.GetEnvOrDefault("CONVERTION_KEY", "conversion")

	vc := converter.NewVideoConverter(db, rabbitmqClient)


	messagesChannel, err := rabbitmqClient.ConsumeMessages(exchangeName, routingKey, queueName)
	if err != nil {
		slog.Error("failed to consume messages", slog.String("error", err.Error()))
	}

	for messageDelivered := range messagesChannel {
		go func (delivery amqp.Delivery)  {
			vc.TaskHandler(delivery)
			
		}(messageDelivered)

	}

}
