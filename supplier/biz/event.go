package biz

import (
	"log"
	"encoding/json"

	"lt-test/supplier/tools"
	. "lt-test/supplier/model"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"lt-test/supplier/mq"
	"strconv"
)

const (
	//事件类型
	UPDATE_EVENT = "update"
	DELETE_EVENT = "delete"
	INSERT_EVENT = "insert"

	//需要监控的表
	TABLE_SKU_SUPPLIER_RELEATION = "sku_supplier_relation"
	TABLE_SKU_SUPPLIER_SYNC      = "sku_supplier_sync"

	//初始mysqldump的行数
	START_UP_SYNC_RECORDS = 1000

	//需要监控binlogFile
	BIN_LOG_FILE = "mysql-bin.000076"
	BIN_LOG_POSITION = 40415958

	//读取的日志文件
	BIN_LOG_FILE_TO_READ = "./binlog.txt"
)

var (
	record           = Record{Rows: make([]SkuSupplierId, 0)}
	skuAndSupplierId SkuSupplierId

	skuAndSupplierIdJson chan []byte = make(chan []byte,2048)
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
	pos := tools.Position{}
	pos,err = tools.ReadFileLast(BIN_LOG_FILE_TO_READ)
	log.Println(pos)
	if err != nil{
		log.Println(err)
	}
	//All(c)
	//return
	if len(pos.FileName) >= 1  && len(pos.Pos) >= 1 {
		Increment(c,pos)
	}else{
		All(c)
	}
	return
}