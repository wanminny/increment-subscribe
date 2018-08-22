package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"lt-test/supplier/retry/db"
	"net/http"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.LstdFlags)
}

var mysqlConfig *db.MySQLClient

type messageProc struct {
	Message  string
	Callback string
	System   string
	Id       int
}

//var (
//	cfg = pflag.StringP("config", "c", "", "config file path.")
//)

func initConfig() {

	//viper.AddConfigPath("./config/")
	//viper.AddConfigPath("/gopath/src/lt-test/supplier/retry/config/")
	//viper.SetConfigType("yml") // 设置配置文件格式为YAML
	viper.SetConfigFile("/gopath/src/lt-test/supplier/retry/config/db.yml")

	//viper.SetConfigFile("db.yml")
	err := viper.ReadInConfig() // 搜索路径，并读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//viper.AutomaticEnv()
	//初始化文件
	env := viper.GetString("db.env")
	host := viper.GetString("mysql." + env + ".host")
	userName := viper.GetString("mysql." + env + ".username")
	password := viper.GetString("mysql." + env + ".password")

	port := viper.GetInt("mysql." + env + ".port")
	dbname := viper.GetString("mysql." + env + ".dbname")

	maxIdle := viper.GetInt("mysql." + env + ".maxIdle")
	maxOpen := viper.GetInt("mysql." + env + ".maxOpen")

	mysqlConfig = &db.MySQLClient{
		Host:    host,
		User:    userName,
		Pwd:     password,
		Port:    port,
		DB:      dbname,
		MaxIdle: maxIdle,
		MaxOpen: maxOpen,
	}

	if mysqlConfig == nil {
		log.Fatal("mysqlconfig is  nil")

	} else {
		err := mysqlConfig.Init()
		if err != nil {
			log.Fatal("mysql init error " + err.Error())
		}
	}

}

func endToOk(m *messageProc) {

	if len(m.Callback) > 0 {
		//
		//http.Get(m.Callback)
		httpGet(m.Id, m.Callback)
		//or post
		//httpPost(m.Id,m.Callback,"")
	}

}

func httpGet(id int, url string) {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err.Error())
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	//log.Println(string(body))

	//成功则 更新状态
	stmt, err := mysqlConfig.Pool.Prepare("update local_tran_message set status = 2 where id =?")

	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	relt, err := stmt.Exec(id)
	if err != nil {
		log.Println(err)
	}
	log.Println(relt.RowsAffected())

}

func httpPost(id int, url string, params string) {

	resp, err := http.Post(url, "application/json", strings.NewReader(params))

	if err != nil {
		log.Fatal(err.Error())
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(string(body))

	//成功则 更新状态
	stmt, err := mysqlConfig.Pool.Prepare("update local_tran_message set status = 2 where id =?")

	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	relt, err := stmt.Exec(id)
	if err != nil {
		log.Println(err)
	}
	log.Println(relt.RowsAffected())

}

func alltryProc() {
	// db schema
	// id , message, create_time,update_time,status,is_del,callback system
	//轮询取出db数据库中的信息;
	//取出没有完成的状态 1 未完成 2 已完成
	// 自己的业务逻辑；
	rows, err := mysqlConfig.Pool.Query("select id, message,callback,system from local_tran_message where status = 1")
	if err != nil {
		log.Fatal("query failed " + err.Error())
	}
	defer rows.Close()

	length := 0
	for rows.Next() {
		length++
		var message string
		var callback string
		var system string
		var id int

		rows.Scan(&id, &message, &callback, &system)
		messProc := &messageProc{
			Id:       id,
			Message:  message,
			Callback: callback,
			System:   system,
		}
		log.Println(length)
		//不断幂等重试 指导成功；可能需要回调;
		go endToOk(messProc)
	}

	// 没有就直接 sleep 1 秒钟
	if length == 0 {
		time.Sleep(1 * time.Second)
	}
}

//将分布式事务转换成本地事务;
func main() {

	//初始化配置
	initConfig()

	//不断重试 幂等 更新状态
	for {
		alltryProc()
	}

}
