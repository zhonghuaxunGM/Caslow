package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type httpError struct {
	Code int
	Mesg string
}

func handler(proc func(url.Values, map[string]interface{}) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var args url.Values
		var out bytes.Buffer
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		requestTime := time.Now()
		defer func() {
			code := http.StatusOK
			data := out.String()
			if e := recover(); e != nil {
				fmt.Errorf("handler defer", "%+v", e)
				switch e.(type) {
				case httpError:
					code = e.(httpError).Code
					data = e.(httpError).Mesg
				default:
					code = http.StatusInternalServerError
					data = e.(error).Error()
				}
				fmt.Println("code:", code)
				if code >= 300 && code < 400 {
					// stringReader := strings.NewReader("test!")
					// r.Body = ioutil.NopCloser(stringReader)
					// fmt.Println("sadsr.Body:", r.Body)
					http.Redirect(w, r, data, code)
				} else {
					http.Error(w, data, code)
				}
			}
			if strings.Contains(args.Get("REQUEST_URL_PATH"), "login") {
				delete(args, "code")
				delete(args, "pass")
			}
			delete(args, "REQUEST_URL_PATH")
			lms <- logMessage{
				Client:   r.RemoteAddr,
				Time:     requestTime,
				Duration: time.Since(requestTime).Seconds(),
				Request:  r.URL.Path,
				Params:   args,
				Cookie:   r.Cookies(),
				Code:     code,
				Reply:    data,
			}
		}()
		// TODO Validate

		r.ParseForm()
		fmt.Println("HTTP Method:", r.Method)
		fmt.Println("HTTP Header:", r.Header)
		fmt.Println("r.URL.Path:", r.URL.Path)

		args = r.Form
		// fmt.Println("args:", args)
		// fmt.Println("get req:", r.FormValue("post1"), r.Form["post1"], r.PostFormValue("post1"))

		// fmt.Println("r.Body:", string(r.Body))
		// buf11 := new(bytes.Buffer)
		// buf11.ReadFrom(r.Body)
		// fmt.Println(buf11.String())

		// defer r.Body.Close()
		reqBody := make(map[string]interface{})
		// if buf11.String() != "" {
		// assert(json.NewDecoder(r.Body).Decode(&reqBody))
		// }
		// 业务层增加map key的判断
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			fmt.Errorf("handler json.NewDecoder", "%s", err.Error())
		}
		fmt.Println("HTTP Body:", reqBody)

		args.Add("REQUEST_URL_PATH", r.URL.Path)
		data := proc(args, reqBody)
		fmt.Println("=====hander done=====")
		if e, ok := data.(httpError); ok {
			panic(httpError{
				Code: http.StatusFound,
				Mesg: fmt.Sprintf("/errpage?err=%s", e.Mesg),
			})
		}
		mw := io.MultiWriter(&out, w)
		enc := json.NewEncoder(mw)
		enc.SetIndent("", "    ")
		assert(enc.Encode(data))
		fmt.Println("=====response done=====")
	}
}
