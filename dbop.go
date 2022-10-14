package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type dsInfo struct {
	Driver string
	Dsn    string
	Name   string
}

var dsns map[string]dsInfo
var dl sync.RWMutex

type DbInfo struct {
	Id     string
	DbUser string
	DbPass string
	DbHost string
	DbPort string
	DbName string
}

func LoadDSNs() {
	dbs := []DbInfo{
		DbInfo{
			Id:     "Id",
			DbUser: "DbUser",
			DbPass: "DbPass",
			DbHost: "Host",
			DbPort: "DbPort",
			DbName: "DbName",
		},
	}
	dl.Lock()
	dsns = make(map[string]dsInfo)
	dl.Unlock()
	go func() {
		for {
			for _, di := range dbs {
				dsn := MySqlDsn(di)
				dl.Lock()
				dsns[di.Id] = dsInfo{Driver: "mysql", Dsn: dsn, Name: di.Id}
				dl.Unlock()
			}
			<-time.After(time.Minute * 2)
		}
	}()
}

func MySqlDsn(di DbInfo) string {
	dsn_tpl := "%s:%s@tcp(%s:%s)/%s"
	return fmt.Sprintf(dsn_tpl, di.DbUser, di.DbPass, di.DbHost,
		di.DbPort, di.DbName)
}

func RangeRows(rows *sql.Rows, proc func()) {
	defer func() {
		if e := recover(); e != nil {
			rows.Close()
			panic(e)
		}
	}()
	for rows.Next() {
		proc()
	}
	assert(rows.Err())
}
