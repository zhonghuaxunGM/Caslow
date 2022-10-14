package main

import (
	"fmt"
	"net/url"
)

type epArg struct {
	Name string
	Hint string
}

type endPoint struct {
	EndPoint string
	Purpose  string
	Note     string
	Args     []epArg
	Body     []epArg
	Returns  []epArg
}

func errpage(args url.Values, reqBody map[string]interface{}) interface{} {
	return map[string]string{
		"ErrorInfo": args.Get("err"),
	}
}

func help(args url.Values, reqBody map[string]interface{}) interface{} {
	fmt.Println("help handle")
	path := args.Get("REQUEST_URL_PATH")
	fmt.Println("args.Get(REQUEST_URL_PATH)", path)
	return []endPoint{
		endPoint{
			EndPoint: "/",
			Purpose:  "NO Setup",
			Note:     "Connect /help FOR More Information",
		},
		endPoint{
			EndPoint: "/errpage",
			Purpose:  "Show Error Info",
			Note:     "Error Details",
			Returns: []epArg{
				epArg{
					Name: "ErrorInfo",
					Hint: "ErrorInfo",
				},
			},
		},
		endPoint{
			EndPoint: "/version",
			Purpose:  "Show Version Info",
			Note:     "Version Details",
			Returns: []epArg{
				epArg{
					Name: "app",
					Hint: "Application Name",
				},
				epArg{
					Name: "version",
					Hint: "Application Version",
				},
			},
		},
		endPoint{
			EndPoint: "/conns",
			Purpose:  "Show Avaliable DB Connections",
			Note:     "更新DB信息需调用/db/进行curl设置",
			Returns: []epArg{
				epArg{
					Name: "driver",
					Hint: "Database Driver",
				},
				epArg{
					Name: "dsn",
					Hint: "Data Source Name",
				},
				epArg{
					Name: "name",
					Hint: "DB Name",
				},
			},
		},
		endPoint{
			EndPoint: "/query",
			Purpose:  "query data",
			Note:     "",
			Args: []epArg{
				epArg{
					Name: "use",
					Hint: "choose db name",
				},
			},
			Body: []epArg{
				epArg{
					Name: "sql",
					Hint: "query SQL",
				},
			},
			Returns: []epArg{
				epArg{
					Name: "columns",
					Hint: "columns values",
				},
				epArg{
					Name: "summary",
					Hint: "Got len(result) row(s) in tqs+tfs (query= tqs; fetch= tfs)",
				},
			},
		},
		endPoint{
			EndPoint: "/exec",
			Purpose:  "exec data",
			Note:     "",
			Args: []epArg{
				epArg{
					Name: "use",
					Hint: "choose db name",
				},
			},
			Body: []epArg{
				epArg{
					Name: "sql",
					Hint: "execute SQL",
				},
			},
			Returns: []epArg{
				epArg{
					Name: "last_insert_id",
					Hint: "last insert id",
				},
				epArg{
					Name: "rows_affected",
					Hint: "rows affected",
				},
				epArg{
					Name: "summary",
					Hint: "Executed statement in elapsed",
				},
			},
		},
	}
}
