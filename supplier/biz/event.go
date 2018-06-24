package biz

import (
	"log"
	"encoding/json"
	"fmt"
	"time"
	"strconv"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"

	"lt-test/supplier/tools"
	. "lt-test/supplier/model"
	"lt-test/supplier/mq"
	"lt-test/supplier/crontab"
	. "lt-test/supplier/env"
	"lt-test/supplier/http"
)

var (
	record           = Record{Rows: make([]SkuSupplierId, 0)}
	skuAndSupplierId SkuSupplierId

	skuAndSupplierIdJson = make(chan []byte,2048)
)

// canal for watch insert update; delete
type MyEventHandler struct {
	canal.DummyEventHandler
}


func (h *MyEventHandler) OnTableChanged(schema string, table string) error {
	log.Println(schema, table)
	return nil
}

func (h *MyEventHandler) OnRotate(roateEvent *replication.RotateEvent) error {
	//log.Printf("%#v", roateEvent)
	return nil
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {

	eventType := e.Action
	tables := e.Table.Name
	recordsLen := len(e.Rows)
	//log.Println(e.Rows)
	if eventType == INSERT_EVENT && recordsLen >= START_UP_SYNC_RECORDS {
		//初始mysqldump ignore;
	} else {
		//strings.Split(tables,".")
		if tables == TABLE_SKU_SUPPLIER_RELEATION {
			//log.Println(e.Action, e.Rows)
			record.Rows = record.Rows[:0:0]
			if eventType == UPDATE_EVENT {
				record.EventType = eventType
				//针对更新每个数组中有前后两项；即更改前；更改后
				var originSupplierId interface{}
				var supplierId interface{}
				originSku := ""
				sku := ""
				for k, v := range e.Rows {
					if k%2 == 0 {
						//originSku = v[1].(string)
						skuAndSupplierId.Id = v[0]
						skuAndSupplierId.Sku = v[1].(string)
						//原始supplierId
						skuAndSupplierId.OriginSupplierId = v[3]
						skuAndSupplierId.CurrentTime = tools.CurrentTime()

						originSku = skuAndSupplierId.Sku
						originSupplierId = skuAndSupplierId.OriginSupplierId
					}
					//变更后的记录
					if k%2 == 1 {
						//新的supp
						sku = v[1].(string)
						//supplierId := v[3]
						skuAndSupplierId.SupplierId = v[3]
						skuAndSupplierId.CurrentTime = tools.CurrentTime()
						supplierId = skuAndSupplierId.SupplierId
						//当且仅当初始 和原来的不一致的时候才变更
						if (sku == originSku) && (supplierId != originSupplierId) {
							record.Rows = append(record.Rows, skuAndSupplierId)
						}
					}
				}
				if len(record.Rows) >= 1 {
					tmpSkuAndSupplierIdJson, _ := json.Marshal(record)
					log.Printf("%s\n", tmpSkuAndSupplierIdJson)
					skuAndSupplierIdJson <- tmpSkuAndSupplierIdJson
					//go mq.Producer(skuAndSupplierIdJson)
				}
			} else if eventType == INSERT_EVENT {
				record.EventType = eventType
				//单个或者批量
				for _, v := range e.Rows {
					skuAndSupplierId.Id = v[0]
					skuAndSupplierId.Sku = v[1].(string)
					skuAndSupplierId.SupplierId = v[3]
					skuAndSupplierId.CurrentTime = tools.CurrentTime()
					record.Rows = append(record.Rows, skuAndSupplierId)
				}
				tmpSkuAndSupplierIdJson, _ := json.Marshal(record)
				log.Printf("%s\n", tmpSkuAndSupplierIdJson)
				skuAndSupplierIdJson <- tmpSkuAndSupplierIdJson
				//go mq.Producer(skuAndSupplierIdJson)

			} else if eventType == DELETE_EVENT {
				record.EventType = eventType
				//单个或者批量
				for _, v := range e.Rows {
					skuAndSupplierId.Id = v[0]
					skuAndSupplierId.Sku = v[1].(string)
					skuAndSupplierId.SupplierId = v[3]
					skuAndSupplierId.CurrentTime = tools.CurrentTime()
					record.Rows = append(record.Rows, skuAndSupplierId)
				}
				tmpSkuAndSupplierIdJson, _ := json.Marshal(record)
				log.Printf("%s\n", tmpSkuAndSupplierIdJson)
				skuAndSupplierIdJson <- tmpSkuAndSupplierIdJson
				//go mq.Producer(skuAndSupplierIdJson)

			} else {
				log.Printf("%s", e.Action)
			}

		} else if tables == TABLE_SKU_SUPPLIER_SYNC {

		}
	}

	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

//全量
func All(c *canal.Canal)  {
	// Start canal
	//从最初的mysqldump 开始；先全量后增量；
	go mq.Producer(skuAndSupplierIdJson)

	c.Run()
}

//增量 	//从指定位置开始增量；show master status
func Increment(c *canal.Canal,pos tools.Position)  {

	go mq.Producer(skuAndSupplierIdJson)

	// 之前是写死
	//binlogFile := BIN_LOG_FILE
	//binlogPos := uint32(BIN_LOG_POSITION)

	binlogFile := pos.FileName
	tempPos,_ := strconv.Atoi(pos.Pos)
	binlogPos := uint32(tempPos)

	p := mysql.Position{
		Name:binlogFile,
		Pos:binlogPos,
	}
	c.RunFrom(p)
}


func Start(c *canal.Canal) (err error) {

	//开启日志查看
	go http.StartHttpService()
	//检查mq是否健康；
	go crontab.CheckMqIsAlive()

	//gtid is ok
	gtidValue := ""
	gOk,err := checkGtid(c)
	if err != nil{
		log.Println(err)
	}
	if gOk {
		gtidValue,err = judgeGtid(c)
		if err != nil{
			log.Println(err)
		}
	}

	//gtid is ok
	if  gOk && len(gtidValue) > 1 {
		posGtid := tools.Position{}
		posGtid,err = tools.ReadFileLast(BIN_LOG_FILE_TO_READ_GTID)
		if err != nil{
			log.Fatal(err)
		}
		log.Println(posGtid.Gtid)

		//定时更新gtid poistion
		go crontUpdateGtidFild(c)

		startFromGtid(c,gtidValue)
	}else {
		pos := tools.Position{}
		pos,err = tools.ReadFileLast(BIN_LOG_FILE_TO_READ)
		if err != nil{
			log.Fatal(err)
		}

		//定时更新 binlog；防止mysql挂掉后重新mysqldump;
		go crontab.ToUpdateBinLogFile()

		// 增量
		if len(pos.FileName) > 1 && len(pos.Pos) >= 1 {
			log.Println("increment")
			Increment(c,pos)
		}else{
			log.Println("all")
			//全量
			All(c)
		}
	}
	return
}

func crontUpdateGtidFild(c *canal.Canal)  {

	timer := time.NewTimer(UPDATE_FILE_IDLE_TIME * time.Hour)
	for {
		timer.Reset(UPDATE_FILE_IDLE_TIME * time.Hour)
		select {
		//case <-time.After(UPDATE_FILE_IDLE_TIME * time.Hour):
		case <-timer.C:
			gtidValue,_ := judgeGtid(c)
			binInfo := fmt.Sprintf("%s,%s\n",tools.CurrentTime(),gtidValue)
			tools.SaveToFile(binInfo,BIN_LOG_FILE_TO_READ_GTID)
		}
	}
}


//是否开启 Gtid
func checkGtid(c *canal.Canal) (b bool, err error) {
	sql := "show variables like '%gtid%'"
	result,err := c.Execute(sql)
	if err != nil{
		log.Fatal(err)
	}
	//log.Printf("%#v",result)
	gtid,err := result.GetString(4,1)

	if gtid != "ON" && gtid != "ON_PERMISSIVE" {
		b = false
	}else{
		b = true
	}
	return
}

// 获取gtid值
func judgeGtid(c *canal.Canal) (gtid string,err error) {
	//sql := "select * from blog2.tbl_comment"
	sql := "show master status"
	result,err := c.Execute(sql)
	if err != nil{
		log.Fatal(err)
	}
	//log.Printf("%#v",result)
	gtid,err = result.GetString(0,4)
	return
}

// gtid mode
func startFromGtid(c *canal.Canal,gtid string)(error)  {

	//pos := c.SyncedPosition()
	//gtid_tmp := c.SyncedGTID()
	//log.Println(pos,gtid_tmp)
	//return nil
	//myGtid := MyGtid{Gtid:col}
	endGtid,err:= mysql.ParseGTIDSet(mysql.MySQLFlavor,gtid)
	if err != nil{
		log.Fatal(err)
	}
	log.Println("gtid mode ")
	c.StartFromGTID(endGtid)
	return nil
}