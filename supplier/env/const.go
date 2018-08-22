package env

const (
	MYSQL_INI_FILE_RELEASE = "mysql_release.ini"
	MYSQL_INI_FILE_TEST    = "mysql_test.ini"

	RABBIT_MQ_FILE_RELEASE = "rabbit_release.ini"
	RABBIT_MQ_FILE_TEST    = "rabbit_test.ini"
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
	BIN_LOG_FILE     = "mysql-bin.000076"
	BIN_LOG_POSITION = 40415958

	//读取的日志文件
	BIN_LOG_FILE_TO_READ = "logs/binlog.txt"

	//读取的GTID日志文件
	BIN_LOG_FILE_TO_READ_GTID = "logs/binlog_gtid.txt"

	//定时刷新gtid binlog file 和普通的binlog file
	UPDATE_FILE_IDLE_TIME = 6
)
