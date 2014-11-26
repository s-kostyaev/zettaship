package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/zazab/zhash"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Message map[string]interface{}

func main() {
	setupLogger()
	message, statusCode, err := sendRequest(createRequest())
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if statusCode != 200 {
		logger.Error(message["error"].(string))
		return
	}
	parseMessage(message)
}

func parseMessage(m Message) {
	hash := zhash.HashFromMap(m)
	format, err := hash.GetString("stdout", "format")
	if err != nil && !zhash.IsNotFound(err) {
		logger.Error(err.Error())
		return
	}

	if format == "table" {
		data, err := hash.GetSlice("stdout", "data")
		if err != nil {
			logger.Error(err.Error())
			return
		}
		table := tablewriter.NewWriter(os.Stdout)
		head := []string{}
		names := []string{}
		for name, _ := range data[0].(map[string]interface{}) {
			head = append(head, strings.ToUpper(name))
			names = append(names, name)
		}
		err = table.Append(head)
		if err != nil {
			logger.Error(err.Error())
		}
		newTable := [][]string{}
		for _, row := range data {
			rowHash := zhash.HashFromMap(row.(map[string]interface{}))
			newRow := []string{}
			for _, name := range names {
				element, err := rowHash.GetString(name)
				if err != nil {
					logger.Error(err.Error())
				}
				newRow = append(newRow, element)
			}
			newTable = append(newTable, newRow)
		}
		err = table.AppendBulk(newTable)
		if err != nil {
			logger.Error(err.Error())
		}

		table.SetBorder(false)
		table.SetRowLine(false)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		table.Render()
		return
	}
	data, err := hash.GetStringSlice("stdout", "data")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, str := range data {
		fmt.Println(str)
	}
	stderr, err := hash.GetStringSlice("stderr")
	if err != nil {
		logger.Error(err.Error())
	}
	for _, str := range stderr {
		fmt.Fprintln(os.Stderr, str)
	}
}

func sendRequest(req *http.Request) (Message, int, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	dec := json.NewDecoder(resp.Body)
	m := Message{}
	if err := dec.Decode(&m); err != nil {
		return nil, 0, err
	}
	return m, resp.StatusCode, nil
}

func createRequest() *http.Request {
	method := ""
	req := &http.Request{}
	err := errors.New("")
	if len(os.Args) == 1 {
		req, err = http.NewRequest("GET", config.ServerUrl, nil)
		if err != nil {
			logger.Error(err.Error())
		}
	} else {
		switch os.Args[1] {
		case "destroy":
			method = "DELETE"
		case "create", "snap", "snapshot", "clone", "set", "rename":
			method = "POST"
		default:
			method = "GET"
		}
		req, err = http.NewRequest(method, config.ServerUrl+os.Args[1], nil)
		if err != nil {
			logger.Error(err.Error())
		}
	}
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
	return req
}
