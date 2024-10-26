package rabbitmq_test

import (
	"context"
	"fmt"
	"gotranscoder/pkg/rabbitmq"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRabbitMQPublishAndConsume(t *testing.T) {
	ctx := context.Background()
	rabbitmqContainer, rabbitmqURL, err := startRabbitMQContainer(ctx)

	assert.NoError(t, err, "Failed to start RabbitMQ testing container")
	defer rabbitmqContainer.Terminate(ctx)

	client, err := rabbitmq.NewRabbitClient(ctx, rabbitmqURL)
	assert.NoError(t, err, "failed to connect to RabbitMQ")
	defer client.Close()

	exchange := "test_exchange"
	routingKey := "test_key"
	queueName := "test_queue"

	msgs, err := client.ConsumeMessages(exchange, routingKey, queueName)
	assert.NoError(t, err, "failed to consume messages")

	t.Run("Publish and consume a valid message", func(t *testing.T) {
		testMessage := []byte("Hello, RabbitMQ!")
		err = client.PublishMessage(exchange, routingKey, queueName, testMessage)
		assert.NoError(t, err, "failed to publish message")

		select {
		case msg := <-msgs:
			assert.Equal(t, string(testMessage), string(msg.Body), "message content mismatched")
			fmt.Println("received message: ", string(msg.Body))
			msg.Ack(false)
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for message")
		}
	})

	t.Run("Consume from a non-existent queue", func(t *testing.T) {
		invalidQueue := "invalid_queue"
		msgs, err := client.ConsumeMessages(exchange, routingKey, invalidQueue)
		assert.NoError(t, err, "consuming from a non-existent queue should not return an error immediately")

		select {
		case <-msgs:
			t.Fatal("Did not expect to receive any messages from non-existent queue")
		case <-time.After(2 * time.Second):
			// pass
		}
	})

	t.Run("Reconnect on connection failure", func(t *testing.T) {
		client.Close()

		err := client.Reconnect(ctx)
		assert.NoError(t, err, "failed to reconnect to RabbitMQ")

		msgs, err = client.ConsumeMessages(exchange, routingKey, queueName)
		assert.NoError(t, err, "failed to consume messages after reconnect")

		err = client.PublishMessage(exchange, routingKey, queueName, []byte("Reconnected message!"))
		assert.NoError(t, err, "failed to publish message after reconnect")

		select {
		case msg := <-msgs:
			assert.Equal(t, "Reconnected message!", string(msg.Body), "message mismatched afted reconnect")
			msg.Ack(false)
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for message after reconnect")
		}
	})
}

func startRabbitMQContainer(ctx context.Context) (rabbitmqContainer testcontainers.Container, url string, err error) {
	req := testcontainers.ContainerRequest{
		Image: "rabbitmq:3-management",
		ExposedPorts: []string{
			"5672/tcp",
			"15672/tcp",
		},
		WaitingFor: wait.ForLog("Server startup complete"),
	}

	rabbitmqContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	host, err := rabbitmqContainer.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	port, err := rabbitmqContainer.MappedPort(ctx, "5672")
	if err != nil {
		return nil, "", err
	}

	rabbitmqURL := fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port.Port())
	fmt.Println("RabbitMQ URL: ", rabbitmqURL)

	return rabbitmqContainer, rabbitmqURL, nil
}
