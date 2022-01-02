package main

import (
	rabbitmq "github.com/ronRabbitMQ/RabbitMQ"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQRounting(
		"routingProduct", "routingKeyOne")
	defer rabbitmq.Destory()

	rabbitmq.RecieveRouting()
}
