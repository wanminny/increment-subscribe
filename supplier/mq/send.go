package mq

// 生产者 简单模式
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"lt-test/supplier/tools"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

//last is vhost;
//var url = "amqp://wanmin:wanmin@localhost:5672/golang"

var (
	Url    = ""
	scheme = "amqp://"
	mqC    = tools.RabbitMqConfig{}

	MAX_RECONNECT_TIME_IDLE = time.Duration(3)
)

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

func Producer(msg chan []byte) {

	mqC.ReadMQIni(&mqC,"mq")
	Url = scheme + mqC.Username + ":" + mqC.Password + "@" + mqC.Host + ":" + mqC.Port + "/" + mqC.Vhost
	queueName := mqC.Queue
	//exchange := mqC.Exchange
	//log.Println(Url, queueName, exchange)
	//拨号；建立连接
	conn, err := amqp.Dial(Url)

	defer conn.Close()
	if err != nil {
		log.Println(err)
		tools.DdTalk([]byte(err.Error()))
	}

	//通过链接建立通道
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		tools.DdTalk([]byte(err.Error()))
	}
	defer ch.Close()

	//可以自己创建；如果不存在；如果已知队列存在的话；可以省略
	//创建（声明）队列
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Println(err)
	}

	for v := range msg {
		//日志
		dumpLog := fmt.Sprintf("%s\n", v)
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
				Body:        v,
			})

		if err != nil {
			log.Println(err)
			tools.DdTalk([]byte(err.Error()))
			ch = reconnect(Url)
			//重新发送
			ch.Publish(
				"",         // exchange
				queue.Name, // routing key
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        v,
				})
		}

	}
}
