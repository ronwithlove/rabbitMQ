package main

import (
	rabbitmq "github.com/ronRabbitMQ/RabbitMQ"
)


func main() {
	rabbitmq:= rabbitmq.NewRabbitMQSimple(
		"ronQueues")
	defer rabbitmq.Destory()

	rabbitmq.ConsumeSiple()
}
