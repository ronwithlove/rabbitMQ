package rabbitmq

import (
	"github.com/streadway/amqp"
)

//创建订阅模式的交换机
func (r *RabbitMQ) exchangeDeclare() {
	err := r.channel.ExchangeDeclare( //如果不存在，创建新的
		r.Exchange,
		"fanout", //交换机类型，这里是广播
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")
}

//订阅模式
func NewRabbitMQPubSub(exchange string) *RabbitMQ {
	//注意，这里不需要设置queue名字，但是exchange还是要的
	return NewRabbitMQ("", exchange, "")
}

//订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
	//1.创建交换机
	r.exchangeDeclare()

	//2.发送消息到队列中
	var err error
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	r.failOnErr(err, "Failed to publish a channel")
}

//订阅模式接受代码
func (r *RabbitMQ) RecieveSub() {
	//1.创建交换机
	r.exchangeDeclare()

	//2.申请队列
	var err error
	q, err := r.channel.QueueDeclare(
		"", //不指定，为空,会自动创建,所以不要指定
		false,
		false,
		true, //排他性开启
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//3.绑定队列到exchange
	err = r.channel.QueueBind(
		q.Name,     //系统自动生成的队列
		"",         //订阅模式key必须为空
		r.Exchange, //我自己命名的交换机
		false,
		nil)
	r.failOnErr(err, "Failed to bind queue")

	//4.接受消息
	r.consume()
}
