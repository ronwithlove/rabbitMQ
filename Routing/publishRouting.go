package main

import (
	"fmt"
	rabbitmq "github.com/ronRabbitMQ/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmqOne := rabbitmq.NewRabbitMQRounting(
		"routingProduct", "routingKeyOne")
	defer rabbitmqOne.Destory()
	rabbitmqTwo := rabbitmq.NewRabbitMQRounting(
		"routingProduct", "routingKeyTwo")
	defer rabbitmqTwo.Destory()

	for i := 0; i <= 100; i++ {
		rabbitmqOne.PublishRouting("Hello World one!" + strconv.Itoa(i))
		rabbitmqTwo.PublishRouting("Hello World two!" + strconv.Itoa(i))
		time.Sleep(time.Second * 1)
		fmt.Println("Hello World!"+strconv.Itoa(i), " 发送成功！")
	}
}
