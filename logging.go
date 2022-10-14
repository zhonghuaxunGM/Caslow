package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type logMessage struct {
	Client   string
	Time     time.Time
	Duration float64
	Request  string
	Params   url.Values
	Cookie   []*http.Cookie
	Code     int
	Reply    string
}

var lms chan logMessage

func init() {
	lms = make(chan logMessage, 1024)
	go func() {
		for {
			lm := <-lms
			lm.Save()
		}
	}()
	if rc.LOG_ROTATE > 0 {
		go func() {
			keepTime := time.Duration(rc.LOG_ROTATE*24) * time.Hour
			for {
				<-time.After(24 * time.Hour)
				lfs, err := filepath.Glob(path.Join(rc.LOG_PATH, "*.log"))
				if err != nil {
					msg := trace("ERROR: %v", err)
					for _, m := range msg {
						fmt.Println(m)
					}
					continue
				}
				for _, lf := range lfs {
					fi, err := os.Stat(lf)
					if err != nil {
						msg := trace("ERROR: %v", err)
						for _, m := range msg {
							fmt.Println(m)
						}
						continue
					}
					if time.Since(fi.ModTime()) > keepTime {
						os.Remove(lf)
					}
				}
			}
		}()
	}
}

func (lm *logMessage) Save() {
	fn := path.Join(rc.LOG_PATH, lm.Time.Format("20060102")+".log")
	logfile, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	assert(err)
	defer logfile.Close()
	_, err = logfile.WriteString(lm.String())
	assert(err)
}

func (lm *logMessage) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "TIME:\t%v\n", lm.Time)
	fmt.Fprintf(&buf, "SPENT:\t%d ms\n", int(lm.Duration*1000+0.5))
	fmt.Fprintf(&buf, "PEER:\t%s\n", lm.Client)
	fmt.Fprintf(&buf, "URI:\t%s\n", lm.Request)
	fmt.Fprintf(&buf, "PARAMS:\n")
	for k, v := range lm.Params {
		fmt.Fprintf(&buf, "\t%s=%s\n", k, v[0])
	}
	fmt.Fprintf(&buf, "COOKIE:\n")
	fmt.Fprintf(&buf, "STATUS:\t%d\n", lm.Code)
	fmt.Fprintf(&buf, "REPLY:\n")
	reply := strings.Split(lm.Reply, "\n")
	for i, rl := range reply {
		fmt.Fprintf(&buf, "\t%s\n", rl)
		if i > 3 && len(reply) > 5 {
			fmt.Fprintf(&buf, "\t...\n")
			break
		}
	}
	fmt.Fprintf(&buf, "========\n")
	return buf.String()
}
