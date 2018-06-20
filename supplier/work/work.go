package main

import (
	"lt-test/supplier/tools"
	"github.com/streadway/amqp"
	"fmt"
	"log"
	"time"
)


// 一个消息只能被一个消费者获取
//
func worker(url string)  {

	fmt.Println(url)
	//建立链接
	conn, err := amqp.Dial(url)
	tools.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//根据链接建立通道
	ch, err := conn.Channel()
	tools.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	//  声明或创建队列（已知queueName队列存在了可以省略此操作);
	q, err := ch.QueueDeclare(
		queueName_tmp, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	tools.FailOnError(err, "Failed to declare a queue")

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
		false,  // auto-ack  //自动确认关闭；比如手动确认
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
			//time.Sleep(time.Second*35)
			time.Sleep(time.Microsecond)
			//log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}



var queueName_tmp = "fuckname"

func main()  {

	//mq URL supplier是指vhost;
	url := "amqp://wanmin:wanmin@127.0.0.1:5672/supplier"

	worker(url)

}