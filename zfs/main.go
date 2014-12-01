package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/zazab/zhash"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Message map[string]interface{}

var (
	commandUrl = config.ServerUrl + "run/"
)

func main() {
	setupLogger()
	message, statusCode, err := sendCommandWithArgs(os.Args)
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

	switch format {
	case "table":
		printTable(hash)
	default:
		printSimple(hash)
	}
}

func printSimple(hash zhash.Hash) {
	data, err := hash.GetStringSlice("stdout", "data")
	if err != nil {
		logger.Error(err.Error())
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

func printTable(hash zhash.Hash) {
	data, err := hash.GetSlice("stdout", "data")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	columnNames, err := hash.GetStringSlice("stdout", "header")
	if err != nil {
		logger.Error(err.Error())
	}
	head := make([]string, len(columnNames))
	for i, name := range columnNames {
		head[i] = strings.ToUpper(name)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Append(head)
	for _, row := range data {
		rowHash := zhash.HashFromMap(row.(map[string]interface{}))
		newRow := make([]string, len(columnNames))
		for i, name := range columnNames {
			element, err := rowHash.GetString(name)
			if err != nil {
				logger.Error(err.Error())
			}
			newRow[i] = element
		}
		table.Append(newRow)
	}

	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
}

func sendCommandWithArgs(args []string) (Message, int, error) {
	uri := ""
	if len(args) == 1 {
		uri = commandUrl
	} else {
		argString := url.QueryEscape(strings.Join(args[1:], " "))
		uri = commandUrl + argString
	}
	resp, err := http.Post(uri, "application/x-www-form-urlencoded", nil)
	if err != nil {
		logger.Error(err.Error())
	}
	dec := json.NewDecoder(resp.Body)
	m := Message{}
	if err := dec.Decode(&m); err != nil {
		return nil, 0, err
	}
	return m, resp.StatusCode, nil
}
