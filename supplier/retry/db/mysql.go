package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

type MySQLClient struct {
	Host    string
	User    string
	Pwd     string
	Port    int
	DB      string
	Pool    *sql.DB
	MaxIdle int
	MaxOpen int
}

func (mc *MySQLClient) Init() (err error) {
	// 构建 DSN 时尤其注意 loc 和 parseTime 正确设置
	// 东八区，允许解析时间字段
	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=%s&parseTime=true",
		mc.User,
		mc.Pwd,
		mc.Host,
		mc.Port,
		mc.DB,
		url.QueryEscape("Asia/Shanghai"))
	// Open 全局一个实例只需调用一次
	mc.Pool, err = sql.Open("mysql", uri)
	if err != nil {
		return err
	}
	//使用前 Ping, 确保 DB 连接正常
	err = mc.Pool.Ping()
	if err != nil {
		return err
	}
	// 设置最大连接数，一定要设置 MaxOpen
	mc.Pool.SetMaxIdleConns(mc.MaxIdle)
	mc.Pool.SetMaxOpenConns(mc.MaxOpen)
	return nil
}
