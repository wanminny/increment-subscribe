package main

// 发布 订阅模式
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
	//queueName := mqC.Queue
	//exchange := mqC.Exchange
	//log.Println(Url,queueName,exchange)
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


	//创建或者声明 交换机 如果交换机不存在就创建；如果已经存在；此操作可以省略;rabbitmq中 交换机没有存储能力
	exchangeName := "fuck_exchange"
	err = ch.ExchangeDeclare(
		exchangeName,   // name  exchange name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)


	length := 100
	for i:=1;i<=length;i++{
		//发送消息
		err = ch.Publish(
			exchangeName,     // exchange
			"", // routing key
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
