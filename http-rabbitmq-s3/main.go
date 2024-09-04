package main

import (
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var RabbitMQURL = "amqp://guest:guest@localhost:5672/"

func main() {
	// move RabbitMQ URL to env variable
	if os.Getenv("RABBITMQ_URL") != "" {
		RabbitMQURL = os.Getenv("RABBITMQ_URL")
	}

	events := NewEvent(RabbitMQURL)

	httpServer := &HttpServer{
		serverConfig: &ServerConfig{
			Event: events,
			host:  "localhost",
			port:  8080,
		},
	}

	events.handleConsumeEvent()

	go httpServer.serveHTTP()

	select {}
}
