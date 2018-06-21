package crontab

import (
	"github.com/siddontang/go-mysql/replication"
	"time"
	"strconv"
	"math/rand"
	"lt-test/supplier/tools"
	"log"
)

var (

	crontabSqlC = &tools.MysqlConfig{}
)

func recordUpdate()  {

	crontabSqlC.ReadMySqlIni(crontabSqlC)
	port,_ := strconv.Atoi(crontabSqlC.Port)


	rand.Seed(time.Now().Unix())
	serverID := uint32(rand.Intn(1000)) + 2001

	// Create a binlog syncer with a unique server id, the server id must be different from other MySQL's.
	// flavor is mysql or mariadb
	cfg := replication.BinlogSyncerConfig {
		ServerID: serverID,
		Flavor:   "mysql",
		Host:     crontabSqlC.Host,
		Port:     uint16(port),
		User:     crontabSqlC.Username,
		Password: crontabSqlC.Password,
	}
	syncer := replication.NewBinlogSyncer(cfg)

	if syncer != nil{
		//注册；便于获取连接；从而定时更新binlog and position;
		bin := syncer.RegisterSlave()
		if bin != nil{
			bin.FlushBinLogRecord()
			//log.Println("to update binlog.txt.....")
		}else{
			log.Println("syncer return err")
		}
	}else{
		log.Println("syncer is nil")
	}
}

//没四个小时刷新一次;防止mysql意外宕机后需要重新同步数据过多
func ToUpdateBinLogFile()  {

	for{
		select {
			//case <-time.After(4*time.Second):
			case <-time.After(2*time.Hour):
				recordUpdate()
		}
	}
}
