package crontab

import (
	"log"

	"github.com/streadway/amqp"
	"lt-test/supplier/tools"
	"time"
)

const (
	PER_TIME_CHECK_IDLE_TIME = 2
)

var (
	mqC    = tools.RabbitMqConfig{}
	scheme = "amqp://"
	Url    = ""
)

// 每两分钟检查一次 是否健康
func CheckMqIsAlive() {

	//time.AfterFunc(PER_TIME_CHECK_IDLE_TIME * time.Minute,isAlive)
	t := time.NewTicker(PER_TIME_CHECK_IDLE_TIME * time.Minute)
	go func() {
		for v := range t.C {
			log.Println(v)
			isAlive()
		}
	}()
}

// 检查rabbit是否健康
func isAlive() {

	mqC.ReadMQIni(&mqC, "mq")
	Url = scheme + mqC.Username + ":" + mqC.Password + "@" + mqC.Host + ":" + mqC.Port + "/" + mqC.Vhost
	queueName := mqC.Queue

	//exchange := mqC.Exchange
	//拨号；建立连接
	conn, err := amqp.Dial(Url)
	if err != nil {
		log.Println(Url, err)
		tools.DdTalk([]byte(err.Error()))
		return
	}
	defer conn.Close()

	//通过链接建立通道
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		tools.DdTalk([]byte(err.Error()))
		return
	}
	defer ch.Close()

	//可以自己创建；如果不存在；如果已知队列存在的话；可以省略
	//创建（声明）队列
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		tools.DdTalk([]byte(err.Error()))
		return
	}

}
