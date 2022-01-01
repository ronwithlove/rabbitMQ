package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const MQURL = "amqp://ron:ron@127.0.0.1:5672/ronMQ" //ronMQ 是virtual host名字

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string //队列
	Exchange  string //交换机
	Key       string //binding key,simple模式用不到
	Mqurl     string
}

//创建RabbitMQ结构体实例
func NewRabbitMQ(queue, exchange, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queue, Exchange: exchange, Key: key, Mqurl: MQURL}
	var err error
	//创建链接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建链接失败")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

//断开channel 和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
	log.Printf("channel, conn断开")

}

//错误处理
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//1.创建simple模式下RabbitMQ实例
func NewRabbitMQSimple(queue string) *RabbitMQ {
	return NewRabbitMQ(queue, "", "") //exchange不写是使用默认的，key不写是没有
}

//2.生产代码
func (r *RabbitMQ) PublishSimple(message string) {
	r.queueDeclare()

	//发送消息到队列中
	err := r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		fmt.Println(err)
	}
}

//申请队列，如果不存在会自动创建
func (r *RabbitMQ) queueDeclare() {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *RabbitMQ) ConsumeSimple() {
	r.queueDeclare()
	r.consume()
}

//接受消息
func (r *RabbitMQ) consume() {
	//接受消息
	msgsChan, err := r.channel.Consume(
		r.QueueName,
		"", //不区分消费者
		true,
		false, //没有排他性
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgsChan {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages, CTRL+C to exit.")

	<-forever
}

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

	//接受消息
	r.consume()
}
