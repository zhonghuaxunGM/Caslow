package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func home(args url.Values, reqBody map[string]interface{}) interface{} {
	fmt.Println("home handle")
	path := args.Get("REQUEST_URL_PATH")
	fmt.Println("args.Get(REQUEST_URL_PATH)", path)
	if path == "/favicon.ico" {
		fmt.Println("favicon.ico")
		return "please add favicon.ico"
	}
	if path == "/" {
		fmt.Println("/")
		return "please access /help"
	}
	return httpError{
		Code: http.StatusNotFound,
		Mesg: "not found: " + path,
	}
}
