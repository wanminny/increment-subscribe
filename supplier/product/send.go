package main

// 生产者 简单模式
import (
	"github.com/streadway/amqp"
	"github.com/astaxie/beego/logs"
	"log"
	"lt-test/supplier/tools"
	"fmt"
	"strconv"
)

func init()  {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

//last is vhost;
//var url = "amqp://wanmin:wanmin@localhost:5672/golang"

var (
	Url = ""
	scheme = "amqp://"
	mqC = tools.RabbitMqConfig{}
)

func Producer()  {

	mqC.ReadMQIni(&mqC)
	Url = scheme + mqC.Username + ":" + mqC.Password +"@" + mqC.Host + ":"+mqC.Port+ "/"+mqC.Vhost
	queueName := mqC.Queue
	exchange := mqC.Exchange

	log.Println(Url,queueName,exchange)
	//拨号；建立连接
	conn,err := amqp.Dial(Url)

	defer conn.Close()
	if err != nil{
		logs.Debug(err)
	}

	//通过链接建立通道
	ch,err := conn.Channel()
	if err != nil{
		logs.Debug(err)
	}
	defer ch.Close()

	//可以自己创建；如果不存在；如果已知队列存在的话；可以省略
	queueName = "fuckname"
	//创建（声明）队列
	queue,err := ch.QueueDeclare(queueName,true,false,false,false,nil)
	if err != nil{
		logs.Debug(err)
	}

	length := 100
	for i:=1;i<=length;i++{
		//发送消息
		err = ch.Publish(
			"",     // exchange
			queue.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(strconv.Itoa(i)),
			})
		if err != nil{
			logs.Debug(err)
		}
	}
}


func main() {
	Producer()
	fmt.Print(Url)
}
