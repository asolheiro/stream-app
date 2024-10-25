package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/streadway/amqp"
)

type RabbitClientInterface interface {
	ConsumeMessages(exchangeName, routingKey, queueName string) (<-chan amqp.Delivery, error)
	PublishMessage(exchangeName, routingKey, queueName string, message []byte) error
	Close() error
	IsClosed() bool
}

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
		conn.Close()
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

func (client *RabbitClient) Close() error {
	if err := client.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}
	if err := client.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	return nil
}

func (client *RabbitClient) IsClosed() bool {
	return client.conn.IsClosed()
}

func (client *RabbitClient) Reconnect(ctx context.Context) (err error) {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled while trying to reconnect")
		default:
			slog.Info("Attempting to reconnect to RabbitMQ...")
			client.conn, client.channel, err = newConnection(client.url)
			if err == nil {
				slog.Info("Reconnected to RabbitMQ successfully")
				return nil
			}
			slog.Error("failed to reconnect to RabbitMQ", slog.String("error", err.Error()))
			time.Sleep(5 * time.Second)
		}
	}
}

func (client *RabbitClient) ConsumeMessages(exchangeName, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	if err := client.channel.ExchangeDeclare(
		exchangeName, "direct",  true, true, false, false, nil,
	); err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	queue, err := client.channel.QueueDeclare(
		queueName, true, true, false, false, nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	if err := client.channel.QueueBind(
		queue.Name, routingKey, exchangeName, false, nil,
	); err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	msgs, err := client.channel.Consume(
		queueName,  "", false, false, false, false, nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %v", err)
	}

	return msgs, nil
}

func (client *RabbitClient) PublishMessage(exchangeName, routingKey, queueName string, message []byte) error {
	if err := client.channel.ExchangeDeclare(
		exchangeName, "direct", true, true, false, false, nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	queue, err := client.channel.QueueDeclare(
		queueName, true, true, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("faile to declare queue: %v", err)
	}

	if err := client.channel.QueueBind(
		queue.Name, routingKey, exchangeName, false, nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	if err := client.channel.Publish(
		exchangeName,  routingKey, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	); err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}
