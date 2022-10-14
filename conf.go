package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func ParseFile(filename string, cfg interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {

		return
	}
	defer f.Close()
	err = ParseReader(f, cfg)
	return
}

func ParseReader(r io.Reader, cfg interface{}) (err error) {
	params := make(map[string]string)
	lines := bufio.NewScanner(r)
	re := regexp.MustCompile(`^[\p{L}_][\p{L}\d_]*$`)
	for lines.Scan() {
		s := lines.Text()
		if strings.HasPrefix(s, "#") || strings.TrimSpace(s) == "" {
			continue
		}
		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			err = fmt.Errorf("non-comment line not in key=value format: %s", s)
			return
		}
		if !re.MatchString(kv[0]) {
			err = fmt.Errorf("key is not valid identifier: [%s]", kv[0])
			return
		}
		params[strings.ToLower(kv[0])] = kv[1]
	}
	err = lines.Err()
	if err != nil {
		return
	}
	err = ParseMap(params, cfg)
	return
}

func ParseMap(m map[string]string, cfg interface{}) (err error) {
	c := reflect.ValueOf(cfg).Elem()
	t := c.Type()
	for i := 0; i < c.NumField(); i++ {
		f := c.Field(i)
		if !f.CanInterface() {
			continue
		}
		key := strings.ToLower(t.Field(i).Name)
		val, ok := m[key]
		if !ok {
			continue
		}
		switch f.Type().Kind() {
		case reflect.String:
			c.Field(i).SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var v int64
			v, err = strconv.ParseInt(val, 10, 64)
			if err != nil {
				return
			}
			c.Field(i).SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:
			var v uint64
			v, err = strconv.ParseUint(val, 10, 64)
			if err != nil {
				return
			}
			c.Field(i).SetUint(v)
		case reflect.Float32, reflect.Float64:
			var v float64
			v, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return
			}
			c.Field(i).SetFloat(v)
		case reflect.Bool:
			var v bool
			v, err = strconv.ParseBool(val)
			if err != nil {
				return
			}
			c.Field(i).SetBool(v)
		}
	}
	return
}

func PersistToMap(cfg interface{}) (m map[string]string, err error) {
	m = make(map[string]string)
	c := reflect.ValueOf(cfg) //.Elem()
	t := c.Type()
	for i := 0; i < c.NumField(); i++ {
		f := c.Field(i)
		if !f.CanInterface() {
			continue
		}
		key := strings.ToLower(t.Field(i).Name)
		switch f.Type().Kind() {
		case reflect.String:
			m[key] = c.Field(i).String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64:
			m[key] = fmt.Sprintf("%v", c.Field(i).Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:
			m[key] = fmt.Sprintf("%v", c.Field(i).Uint())
		case reflect.Float32, reflect.Float64:
			m[key] = fmt.Sprintf("%v", c.Field(i).Float())
		case reflect.Bool:
			m[key] = fmt.Sprintf("%v", c.Field(i).Bool())
		}
	}
	return
}

func PersistToWriter(cfg interface{}, w io.Writer) (err error) {
	m, err := PersistToMap(cfg)
	if err != nil {
		return
	}
	for k, v := range m {
		_, err = w.Write([]byte(fmt.Sprintf("%s=%s\n", k, v)))
		if err != nil {
			return
		}
	}
	return
}

func PersistToFile(cfg interface{}, filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	err = PersistToWriter(cfg, f)
	return
}
