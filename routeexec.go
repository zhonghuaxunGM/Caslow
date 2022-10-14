package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func doexe(conn *sql.DB, args url.Values, reqBody map[string]interface{}) (queryResults, float64) {
	// qry := args.Get("sql")
	qry, ok := reqBody["sql"]
	if !ok {
		return queryResults{}, -1
	}
	ctx, cf := context.WithTimeout(context.Background(), time.Duration(rc.EXEC_TIMEOUT)*time.Second)
	defer cf()
	start := time.Now()
	qr, err := conn.ExecContext(ctx, qry.(string))
	elapsed := time.Since(start).Seconds()
	assert(err)
	var LastInsertId, RowsAffected string
	lid, err := qr.LastInsertId()
	if err == nil {
		LastInsertId = fmt.Sprintf("%d", lid)
	} else {
		LastInsertId = err.Error()
	}
	ra, err := qr.RowsAffected()
	if err == nil {
		RowsAffected = fmt.Sprintf("%d", ra)
	} else {
		RowsAffected = err.Error()
	}
	return queryResults{
		queryResult{
			"last_insert_id": LastInsertId,
			"rows_affected":  RowsAffected,
		},
	}, elapsed
}

func execute(args url.Values, reqBody map[string]interface{}) (res interface{}) {
	use := args.Get("use")
	// qry := args.Get("sql")
	// if use == "" || qry == "" {
	// 	return httpError{
	// 		Code: http.StatusSeeOther,
	// 		Mesg: "/uisql?action=exec&use=" + use,
	// 	}
	// }
	dl.RLock()
	ds, ok := dsns[use]
	dl.RUnlock()
	if !ok {
		return httpError{
			Code: http.StatusNotFound,
			Mesg: "[use] is not a valid data source",
		}
	}
	var (
		dss []dsInfo
		els float64
		rec queryResults
	)
	dss = append(dss, ds)
	defer func() {
		if e := recover(); e != nil {
			fmt.Errorf("execute sql", "%s", e.(error).Error())
			res = httpError{
				Code: http.StatusInternalServerError,
				Mesg: e.(error).Error(),
			}
		}
	}()
	for _, ds := range dss {
		conn, err := getDB(ds.Driver, ds.Dsn, "exec")
		assert(err)
		data, elapsed := doexe(conn, args, reqBody)
		if elapsed == -1 {
			assert(fmt.Errorf("request body lose SQL"))
		}
		els += elapsed
		for _, d := range data {
			if len(dss) > 1 {
				d[rc.DB_TAG] = ds.Name
			}
			rec = append(rec, d)
		}
	}
	summary := fmt.Sprintf("Executed statement in %fs", els)
	rec = append(rec, map[string]interface{}{
		"summary": summary,
	})
	// args.Set("RESTIQUE_SUMMARY", summary)
	return rec
}
