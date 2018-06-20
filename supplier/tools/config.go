package tools

import (
	"github.com/astaxie/beego/config"
	"lt-test/supplier/env"
	"log"
)

// mysql struct
type MysqlConfig struct {
	Host string
	Port string
	Username string
	Password string
	Database string
	Tables string
}

// rabbit mq  struct
type RabbitMqConfig struct {
	Host string
	Port string
	Username string
	Password string

	Vhost string
	Exchange string
	Queue string
}

// 获取mysql ini 配置
func (sql MysqlConfig)ReadMySqlIni(result *MysqlConfig)  {

	ini,err := config.NewConfig("ini","./supplier/config/" + env.MYSQL_INI_FILE_TEST)
	FailOnError(err,"failed to new config")

	result.Host = ini.String("mysql::host")
	result.Port = ini.String("mysql::port")
	result.Username = ini.String("mysql::username")
	result.Password = ini.String("mysql::password")

	result.Database = ini.String("mysql::database")
	result.Tables = ini.String("mysql::tables")
}

// 获取mq ini 配置
func (mq RabbitMqConfig)ReadMQIni(result *RabbitMqConfig)  {

	ini,err := config.NewConfig("ini","./supplier/config/"+ env.RABBIT_MQ_FILE_TEST)
	//FailOnError(err,err.Error())
	if err != nil{
		log.Fatal(err)
	}

	result.Host = ini.String("mq::host")
	result.Port = ini.String("mq::port")
	result.Username = ini.String("mq::username")
	result.Password = ini.String("mq::password")
	result.Vhost = ini.String("mq::vhost")
	result.Exchange = ini.String("mq::exchange")
	result.Queue = ini.String("mq::queue")
}