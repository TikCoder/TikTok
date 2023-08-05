package db

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// 假设有 4 个连接
// 1. DEFAULT_MAX_CONN = 5, DEFAULT_MAX_IDLE_CONN = 3
// 此时可以 同时 查询完毕 4 个连接，但是 连接池 里，能保持 3 个连接等待，剩下一个断开
// 2. DEFAULT_MAX_CONN = 3, DEFAULT_MAX_IDLE_CONN = 3
// 此时 不可以 同时 查询完毕 4 个连接，最多只能开 3 个
// https://juejin.cn/post/6844904087427776519
//const (
//	DEFAULT_MAX_CONN      = 0  // 与数据库连接的最大连接数(允许的最大连接数)
//	DEFAULT_MAX_IDLE_CONN = 0  // 空闲池的最大连接数(连接池里的最大连接数)
//	DEFAULT_IDLE_TIMEOUT  = 10 // 设置最大空闲时间
//	DEFAULT_READ_TIMEOUT  = 10
//	DEFAULT_WRITE_TIMEOUT = 10
//)

// GormLogger Gorm用来打日志的结构体
type GormLogger struct {
}

//var (
//	Factory GormFactory
//)

//// GormFactory 用来生成Gorm指针的工厂
//type GormFactory struct {
//	MaxIdleConn  int
//	MaxConn      int
//	IdleTimeout  int
//	ReadTimeout  int
//	WriteTimeout int
//}

// CreateGorm 创建 sql 连接，并设置最大连接数
func CreateGorm(dsn string, maxConn int, maxIdleConn int, idleTimeout int) (*gorm.DB, error) {
	//// 读超时
	//if p.ReadTimeout == 0 {
	//	p.ReadTimeout = DEFAULT_READ_TIMEOUT
	//}
	//// 写超时
	//if p.WriteTimeout == 0 {
	//	p.WriteTimeout = DEFAULT_WRITE_TIMEOUT
	//}
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return db, err
	}
	logger := &GormLogger{}
	db.LogMode(true)
	db.SetLogger(logger)
	//maxIdleConn := p.MaxIdleConn
	//if maxIdleConn == 0 {
	//	maxIdleConn = DEFAULT_MAX_IDLE_CONN
	//}
	//
	//maxConn := p.MaxConn
	//if maxConn == 0 {
	//	maxConn = DEFAULT_MAX_CONN
	//}
	//idleTimeout := p.IdleTimeout
	//if idleTimeout == 0 {
	//	idleTimeout = DEFAULT_IDLE_TIMEOUT
	//}
	db.DB().SetMaxIdleConns(maxIdleConn)
	db.DB().SetMaxOpenConns(maxConn)
	db.DB().SetConnMaxLifetime(time.Duration(idleTimeout) * time.Second)
	return db, err
}

// Print 实现logger的print函数
func (logger *GormLogger) Print(values ...interface{}) {
	var (
		level = values[0]
	)
	if level == "sql" {
		log.Printf("%+v %s \"\"", values, level)
	} else {
		log.Printf("%+v", values)
	}
}
