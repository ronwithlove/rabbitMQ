package main

import (
	"fmt"
	rabbitmq "github.com/ronRabbitMQ/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQPubSub(
		"newProduct")
	defer rabbitmq.Destory()

	for i := 0; i <= 100; i++ {
		rabbitmq.PublishPub("Hello World!" + strconv.Itoa(i))
		time.Sleep(time.Second * 1)
		fmt.Println("Hello World!"+strconv.Itoa(i), " 发送成功！")
	}
}
