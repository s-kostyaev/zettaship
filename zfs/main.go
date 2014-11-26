package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/zazab/zhash"
	"net/http"
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
		names, err := hash.GetStringSlice("stdout", "header")
		if err != nil {
			logger.Error(err.Error())
		}
		for _, name := range names {
			head = append(head, strings.ToUpper(name))
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
	if len(os.Args) == 1 {
		req, err := http.NewRequest("POST", config.ServerUrl+"run/", nil)
		if err != nil {
			logger.Error(err.Error())
		}
		return req
	}
	args := strings.Join(os.Args[2:], "+")
	req, err := http.NewRequest("POST", config.ServerUrl+"run/"+os.Args[1]+"/"+
		args, nil)
	if err != nil {
		logger.Error(err.Error())
	}

	return req
}
