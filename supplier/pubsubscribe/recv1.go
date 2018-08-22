package main

// 订阅
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"lt-test/supplier/tools"
	"time"
)

func Consumer(url string) {

	fmt.Println(url)
	//建立链接
	conn, err := amqp.Dial(url)
	tools.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//根据链接建立通道
	ch, err := conn.Channel()
	tools.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	//创建或者声明 交换机 如果交换机不存在就创建；如果已经存在；此操作可以省略
	//注意订阅模式是向交换机发送消息；而简单，工作模式是向队列发送消息；【重要！】
	//rabbitmq中 交换机没有存储能力

	exchangeName := "fuck_exchange"
	err = ch.ExchangeDeclare(
		exchangeName, // name  exchange name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	tools.FailOnError(err, "failed to declare exchange!")

	//  声明或创建队列（已知queueName队列存在了可以省略此操作);
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	tools.FailOnError(err, "Failed to declare a queue")

	//将队列绑定到交换机上
	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil)
	tools.FailOnError(err, "Failed to bind queue")

	// 打破平均分配；默认是平均分配
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	tools.FailOnError(err, "Failed to set QoS")

	//消费
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack  //自动确认关闭；必须手动确认
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	tools.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			//故意超时
			time.Sleep(10000 * time.Microsecond)
			//time.Sleep(time.Second*35)
			//log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//订阅模式 fanout 是多个队列 即同一份源发送到多个目标队列；异构数据
var queueName = "fuck_exchange_queue1"

func main() {
	//mq URL supplier是指vhost;
	url := "amqp://wanmin:wanmin@127.0.0.1:5672/supplier"
	Consumer(url)
}
