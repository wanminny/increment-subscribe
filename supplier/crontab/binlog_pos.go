package crontab

import (
	"github.com/siddontang/go-mysql/replication"
	"log"
	. "lt-test/supplier/env"
	"lt-test/supplier/tools"
	"math/rand"
	"strconv"
	"time"
)

var (
	crontabSqlC = &tools.MysqlConfig{}
)

//
func recordUpdate() {

	crontabSqlC.ReadMySqlIni(crontabSqlC)
	port, _ := strconv.Atoi(crontabSqlC.Port)

	rand.Seed(time.Now().Unix())
	serverID := uint32(rand.Intn(1000)) + 2001

	// Create a binlog syncer with a unique server id, the server id must be different from other MySQL's.
	// flavor is mysql or mariadb
	cfg := replication.BinlogSyncerConfig{
		ServerID: serverID,
		Flavor:   "mysql",
		Host:     crontabSqlC.Host,
		Port:     uint16(port),
		User:     crontabSqlC.Username,
		Password: crontabSqlC.Password,
	}
	syncer := replication.NewBinlogSyncer(cfg)

	if syncer != nil {
		//注册；便于获取连接；从而定时更新binlog and position;
		bin := syncer.RegisterSlave()
		if bin != nil {
			bin.FlushBinLogRecord()
			//log.Println("to update binlog.txt.....")
		} else {
			log.Println("syncer return err")
		}
	} else {
		log.Println("syncer is nil")
	}
}

//每6个小时刷新一次;防止mysql意外宕机后需要重新同步数据过多
func ToUpdateBinLogFile() {

	timer := time.NewTimer(UPDATE_FILE_IDLE_TIME * time.Hour)
	for {
		timer.Reset(UPDATE_FILE_IDLE_TIME * time.Hour)
		select {
		//case <-time.After(4*time.Second):
		case <-timer.C: //比after方式节省timer 资源
			recordUpdate()
		}
	}
}
