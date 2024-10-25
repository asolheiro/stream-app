package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return conn, channel, nil
}

func NewRabbitClient(connectionURL string) (*RabbitClient, error) {

	conn, channel, err := newConnection(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return &RabbitClient{
		conn:    conn,
		channel: channel,
		url:     connectionURL,
	}, nil
}

func (client *RabbitClient) Close() {
	client.channel.Close()
	client.conn.Close()
}

func (client *RabbitClient) ConsumeMessages(exchangeName, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	err := client.channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	queue, err := client.channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	err = client.channel.QueueBind(
		queue.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	msgs, err := client.channel.Consume(
		queueName,
		"go-transcoder",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("error")
	}

	return msgs, nil
}

func (client *RabbitClient) PublishMessage(exchangeName, routingKey, queueName string, message []byte) error {
	err := client.channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	queue, err := client.channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("faile to declare queue: %v", err)
	}

	err = client.channel.QueueBind(
		queue.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to bin queue: %v", err)
	}

	err = client.channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}
