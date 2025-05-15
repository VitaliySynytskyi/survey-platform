package connection_examples

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQConfig holds the connection details
type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
}

// ConnectToRabbitMQ establishes a connection to RabbitMQ
func ConnectToRabbitMQ(config RabbitMQConfig) (*amqp.Connection, *amqp.Channel, error) {
	// Create connection URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%s",
		config.User, config.Password, config.Host, config.Port)

	// Add virtual host if specified
	if config.VHost != "" && config.VHost != "/" {
		url = fmt.Sprintf("%s/%s", url, config.VHost)
	}

	// Connect to RabbitMQ server
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	log.Println("Successfully connected to RabbitMQ")
	return conn, ch, nil
}

// Example usage:
func RabbitMQExample() {
	config := RabbitMQConfig{
		Host:     "localhost", // or "rabbitmq" when using Docker network
		Port:     "5672",
		User:     "rabbit_user",
		Password: "rabbit_password",
		VHost:    "/",
	}

	conn, ch, err := ConnectToRabbitMQ(config)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	// Example of declaring a queue
	queueName := "example_queue"
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Example of publishing a message
	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf("Sent %s to queue %s", body, queueName)

	// Example of consuming messages (commented out as it would block indefinitely)
	/*
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Fatalf("Failed to register a consumer: %v", err)
		}

		forever := make(chan bool)
		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
			}
		}()

		log.Printf("Waiting for messages. To exit press CTRL+C")
		<-forever
	*/
}
