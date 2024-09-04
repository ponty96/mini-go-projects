package main

import (
	"context"
	"log"

	// "sync"
	"time"

	"github.com/ponty96/my-proto-schemas/output/schemas"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type Events struct {
	// consider a different connect for publisher and consumer
	conn *amqp.Connection
}

func NewEvent(RabbitMQURL string) Events {
	conn, err := amqp.Dial(RabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	return Events{
		conn: conn,
	}
}

func (e *Events) close() {
	e.conn.Close()
}

func (e *Events) handleConsumeEvent() {
	ch, err := e.conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"test-profiles", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		"profiles", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			// d.Body is a []byte, and we will protobuf encode it
			log.Printf("Received a message: %s", d.Body)
			profileEvent := &schemas.Profile{}

			if err := proto.Unmarshal(d.Body, profileEvent); err != nil {
				// return fmt.Errorf("failed to parse message: %v", err)
				// fmt.Println(err)
				log.Panicf("failed to decode %s", err)
			}
			firstName := profileEvent.GetFirstname()
			lastName := profileEvent.GetLastname()
			log.Printf("Received a message: %s %s", firstName, lastName)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}

func (e *Events) handlePublishEvent(profileJson ProfileJson) {
	ch, err := e.conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"profiles", // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileSchema := &schemas.Profile{
		Firstname: profileJson.FirstName,
		Lastname:  profileJson.LastName,
		Email:     profileJson.Email,
		Phone:     profileJson.Phone,
		Address: &schemas.Address{
			Street:  profileJson.Address.Street,
			City:    profileJson.Address.City,
			State:   profileJson.Address.State,
			Zip:     profileJson.Address.ZipCode,
			Country: profileJson.Address.Country,
		},
		Company: "ACME LTD",
	}

	profilebytes, err := proto.Marshal(profileSchema)

	if err != nil {
		log.Panicf("failed to encode %s", err)
	}
	err = ch.PublishWithContext(ctx,
		"profiles", // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(profilebytes),
		})
	failOnError(err, "Failed to publish a message")
	// log.Printf(" [x] Sent %s", profilebytes)
}
