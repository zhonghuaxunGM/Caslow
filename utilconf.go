package main

import (
	"os"
)

type restiqueConf struct {
	AUTH_TOKEN   string
	SERVICE_PORT string
	LOG_PATH     string
	LOG_ROTATE   int
	PID_FILE     string
	// CLIENT_CIDRS  string
	// HIST_PATH string
	// OPEN_HATEOAS  bool
	QUERY_TIMEOUT int
	EXEC_TIMEOUT  int
	DB_TAG        string
	QUERY_MAXROWS int
}

type DBConf struct {
	IDLE_TIMEOUT int
	SESSION_LIFE int
	// QUERY_MAXROWS int
	// DB_TAG        string
	PFS_HOST string
	PFS_PORT string
}

var rc restiqueConf

func parseConfig(fn string) {
	rc.AUTH_TOKEN = "anything"
	rc.SERVICE_PORT = "32779"
	// rc.IDLE_TIMEOUT = 300
	// rc.SESSION_LIFE = 3600
	rc.DB_TAG = "[DB]"
	rc.PID_FILE = "./caslow.pid"
	rc.LOG_PATH = "./logs"
	// rc.HIST_PATH = "./history"
	// rc.PFS_HOST = "pfs.paadoo.net"
	// rc.PFS_PORT = "2000"
	if fn != "" {
		assert(ParseFile(fn, &rc))
	}
	// if rc.IDLE_TIMEOUT > 86400 {
	// 	rc.IDLE_TIMEOUT = 86400
	// }
	// if rc.SESSION_LIFE > 864000 {
	// 	rc.SESSION_LIFE = 864000
	// }
	if rc.QUERY_TIMEOUT <= 0 {
		rc.QUERY_TIMEOUT = 60
	}
	if rc.EXEC_TIMEOUT <= 0 {
		rc.EXEC_TIMEOUT = 60
	}
	// rc.CLIENT_CIDRS = strings.TrimSpace(rc.CLIENT_CIDRS)
	// if len(rc.CLIENT_CIDRS) > 0 {
	// 	allowed_cidrs = strings.Split(rc.CLIENT_CIDRS, ",")
	// }
	assert(os.MkdirAll(rc.LOG_PATH, 0755))
	// assert(os.MkdirAll(rc.HIST_PATH, 0755))
}
