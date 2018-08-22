package main

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/streadway/amqp"
	"log"
	"lt-test/supplier/tools"
	"time"
)

var (
	scheme = "amqp://"

	//源
	destUrl = ""
	destMqC = tools.RabbitMqConfig{}

	//目的地
	sourceUrl = ""
	sourceMqC = tools.RabbitMqConfig{}

	v = make(chan []byte, 512)
)

const MAX_RECONNECT_TIME_IDLE = 3

// mq 跨机房同步

func Consumer() {

	sourceMqC.ReadMQIni(&sourceMqC, "mq_source")
	sourceUrl = scheme + sourceMqC.Username + ":" + sourceMqC.Password + "@" + sourceMqC.Host + ":" + sourceMqC.Port + "/" + sourceMqC.Vhost
	queueName := sourceMqC.Queue

	//建立链接
	conn, err := amqp.Dial(sourceUrl)
	tools.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//根据链接建立通道
	ch, err := conn.Channel()
	tools.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

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

	// 打破平均分配；默认是平均分配
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	tools.FailOnError(err, "Failed to set QoS")

	forever := make(chan bool)
	go func() {
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
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			//fmt.Printf("%#v",d.Body)

			v <- d.Body
			go Producer(v)

			d.Ack(false)
			//if err == nil{
			//}else{
			//	log.Println(err)
			//}
			//故意超时
			//time.Sleep(1000 * time.Microsecond)
			//time.Sleep(time.Second*35)
			//log.Printf("Done")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func reconnect(Url string) (availableCh *amqp.Channel) {
	//reconnect
	for {
		conn, err := amqp.Dial(Url)
		if err != nil {
			// 钉钉报警
			log.Println(err)
			tools.DdTalk([]byte(err.Error()))
			time.Sleep(MAX_RECONNECT_TIME_IDLE * time.Second)
		} else {
			ch, err := conn.Channel()
			if err != nil {
				log.Println(err)
				tools.DdTalk([]byte(err.Error()))
				time.Sleep(MAX_RECONNECT_TIME_IDLE * time.Second)
			} else {
				//ok
				availableCh = ch
				return
			}
		}
	}
}

func Producer(msg chan []byte) (err error) {

	destMqC.ReadMQIni(&destMqC, "mq_destination")
	destUrl = scheme + destMqC.Username + ":" + destMqC.Password + "@" + destMqC.Host + ":" + destMqC.Port + "/" + destMqC.Vhost
	queueName := destMqC.Queue
	//exchange := mqC.Exchange
	//log.Println(Url, queueName, exchange)
	//拨号；建立连接
	conn, err := amqp.Dial(destUrl)

	defer conn.Close()
	if err != nil {
		log.Println(err)
		tools.DdTalk([]byte(err.Error()))
		return errors.New(err.Error())
	}

	//通过链接建立通道
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		tools.DdTalk([]byte(err.Error()))
		return errors.New(err.Error())

	}
	defer ch.Close()

	//可以自己创建；如果不存在；如果已知队列存在的话；可以省略
	//创建（声明）队列
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return errors.New(err.Error())

	}

	for value := range msg {
		//日志
		dumpLog := fmt.Sprintf("%s\n", value)
		tools.LogTofile(dumpLog)
		//  mysqldump大量记录测试和 时间先后顺序测试；
		//time.Sleep(10000*time.Microsecond)
		//发送消息
		err = ch.Publish(
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        value,
			})

		if err != nil {
			log.Println(err)
			tools.DdTalk([]byte(err.Error()))
			ch = reconnect(destUrl)
			//重新发送
			ch.Publish(
				"",         // exchange
				queue.Name, // routing key
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        value,
				})
		}
	}
	return nil
}

func main() {

	Consumer()

	////go Producer()
	//select {
	//
	//}
}
