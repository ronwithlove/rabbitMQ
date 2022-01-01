package main

import (
	"fmt"
	rabbitmq "github.com/ronRabbitMQ/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQSimple(
		"ronQueues")
	defer rabbitmq.Destory()

	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("Hello World!"+strconv.Itoa(i))
		time.Sleep(time.Second * 1)
		fmt.Println("发送成功！")
	}
}
