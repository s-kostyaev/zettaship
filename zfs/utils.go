package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crackcomm/go-clitable"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Message map[string]interface{}

var (
	client = &http.Client{}
	req    = &http.Request{}
	err    = errors.New("")
	method = ""
)

func parseMessage(m Message, statusCode int) {
	if statusCode != 200 {
		logger.Error(m["error"].(string))
		return
	}
	data := m["stdout"].(map[string]interface{})["data"].([]interface{})
	format, ok := m["stdout"].(map[string]interface{})["format"].(string)
	if ok && format == "table" {
		head := []string{}
		for name, _ := range data[0].(map[string]interface{}) {
			head = append(head, name)
		}
		tables := []map[string]interface{}{}
		for _, table := range data {
			tables = append(tables, table.(map[string]interface{}))
		}
		clitable.PrintTable(head, tables)
		return
	}
	for _, str := range data {
		fmt.Println(str.(string))
	}
	for _, str := range m["stderr"].([]interface{}) {
		fmt.Fprintln(os.Stderr, str.(string))
	}
}

func sendRequest() (Message, int) {
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
	}
	dec := json.NewDecoder(resp.Body)
	m := Message{}
	if err := dec.Decode(&m); err != nil {
		logger.Error(err.Error())
	}
	return m, resp.StatusCode
}

func xorBytes(b1 [32]byte, bmore ...[32]byte) []byte {
	rv := make([]byte, len(b1))

	for i := range b1 {
		rv[i] = b1[i]
		for _, m := range bmore {
			rv[i] = rv[i] ^ m[i]
		}
	}

	return rv
}

func createRequest() {
	if len(os.Args) == 1 {
		req, err = http.NewRequest("GET", config.ServerUrl, nil)
		if err != nil {
			logger.Error(err.Error())
		}
	} else {
		switch os.Args[1] {
		case "destroy":
			method = "DELETE"
		case "mount":
			method = "POST"
			if len(os.Args) == 2 {
				method = "GET"
			}
		case "umount", "unmount", "create", "snap", "snapshot":
			method = "POST"
		default:
			method = "GET"
		}
		req, err = http.NewRequest(method, config.ServerUrl+os.Args[1], nil)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func addArgs() {
	if len(os.Args) > 2 {
		limit := len(os.Args)
		params := url.Values{}
		if strings.Contains(os.Args[len(os.Args)-1], "/") {
			params.Add("last", os.Args[len(os.Args)-1])
			limit = limit - 1
		}
		if len(os.Args) > 3 {
			prev := ""
			for _, arg := range os.Args[2:limit] {
				if prev == "" {
					prev = arg
					continue
				}
				if strings.HasPrefix(arg, "-") {
					params.Add(prev, "")
					prev = arg
					continue
				}
				params.Add(prev, arg)
				prev = ""
			}
			if prev != "" {
				params.Add(prev, "")
			}
		}
		req.URL.RawQuery = params.Encode()
	}
}

func addCookie() {
	cookie := &http.Cookie{}
	cookie.Name = config.Container
	cookie.Value = hex.EncodeToString(xorBytes(sha256.Sum256([]byte(
		config.Container)), sha256.Sum256([]byte(config.Salt))))
	req.AddCookie(cookie)
}
