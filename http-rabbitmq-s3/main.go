package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	AMQP_URL := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(AMQP_URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"events.profiles", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)

	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // name
		"",     // consumer
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)

	go handleMessages(msgs)

	var forever = make(chan bool)

	<-forever
}

func handleMessages(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		fmt.Println(d)
		log.Printf("Received a message: %s", d.Body)
	}
}
