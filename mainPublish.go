package main

import (
	"fmt"
	rabbitmq "github.com/ronRabbitMQ/RabbitMQ"
)


func main() {
	rabbitmq:= rabbitmq.NewRabbitMQSimple(
		"ronQueues")
	defer rabbitmq.Destory()

	rabbitmq.PublishSimple("Hello World!")
	fmt.Println("发送成功！")
}
