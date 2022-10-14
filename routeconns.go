package main

import (
	"database/sql"
	"net/url"
	"sync"
	"time"
)

var connsCache map[string]*sql.DB

func init() {
	connsCache = make(map[string]*sql.DB)
}

func conns(url.Values, map[string]interface{}) interface{} {
	dl.RLock()
	defer dl.RUnlock()
	var cs []map[string]string
	for k, v := range dsns {
		info := map[string]string{
			"name":   k,
			"driver": v.Driver,
			"dsn":    v.Dsn,
		}
		cs = append(cs, info)
	}
	return cs
}

var connsLock sync.RWMutex

func getDB(driver, dsn, mode string) (db *sql.DB, err error) {
	defer func() {
		if err != nil {
			// if mode
		}
	}()
	connsLock.RLock()
	db, ok := connsCache[driver+"/"+dsn]
	connsLock.RUnlock()
	if ok {
		return db, nil
	}
	db, err = sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Second)
	db.SetMaxOpenConns(5)
	connsLock.Lock()
	connsCache[driver+"/"+dsn] = db
	connsLock.Unlock()
	return db, nil
}
