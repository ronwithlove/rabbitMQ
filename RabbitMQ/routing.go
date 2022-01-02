package rabbitmq

import (
	"github.com/streadway/amqp"
)

//路由模式
func NewRabbitMQRounting(exchange, routingKey string) *RabbitMQ {
	//注意 key 不能为空
	return NewRabbitMQ("", exchange, routingKey)
}

//创建路由模式的交换机
func (r *RabbitMQ) routingExchangeDeclare() {
	err := r.channel.ExchangeDeclare( //如果不存在，创建新的
		r.Exchange,
		"direct", //交换机类型，这里是广播
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")
}

//发送
func (r *RabbitMQ) PublishRouting(message string) {
	//1.声明交换机
	r.routingExchangeDeclare()

	//2.发送消息到队列中
	var err error
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	r.failOnErr(err, "Failed to publish a channel")
}

//接收
func (r *RabbitMQ) RecieveRouting() {
	//1.声明交换机
	r.routingExchangeDeclare()

	//2.申请队列，和订阅模式一样，不需要指定队列
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
		r.Key,      //需要设置绑定key
		r.Exchange, //我自己命名的交换机
		false,
		nil)
	r.failOnErr(err, "Failed to bind queue")

	//4.接受消息
	r.consume()
}
