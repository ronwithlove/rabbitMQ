package main

import (
	rabbitmq "github.com/ronRabbitMQ/RabbitMQ"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQPubSub(
		"newProduct")
	defer rabbitmq.Destory()

	rabbitmq.RecieveSub()
}
