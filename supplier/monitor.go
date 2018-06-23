package main

import (
	"strings"
	"log"

	"github.com/astaxie/beego/logs"
	"github.com/siddontang/go-mysql/canal"

	"lt-test/supplier/tools"
	"lt-test/supplier/biz"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate)
}

var(
	sqlC = &tools.MysqlConfig{}
)

func main() {

	sqlC.ReadMySqlIni(sqlC)

	cfg := canal.NewDefaultConfig()
	cfg.Addr = sqlC.Host + ":" + sqlC.Port
	cfg.User = sqlC.Username
	cfg.Password = sqlC.Password
	cfg.Dump.TableDB = sqlC.Database
	tables := strings.Split(sqlC.Tables, ",")
	cfg.Dump.Tables = tables

	c, err := canal.NewCanal(cfg)
	if err != nil {
		logs.Debug(err)
	}

	// 注册handler处理RowsEvent
	c.SetEventHandler(&biz.MyEventHandler{})

	biz.Start(c)
}
