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