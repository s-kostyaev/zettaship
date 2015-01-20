package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/zazab/zhash"
)

type Reply map[string]interface{}

var (
	commandUrl = config.ServerUrl + "run/"
)

func main() {
	setupLogger()
	reply, statusCode, err := sendCommandWithArgs(os.Args)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if statusCode != 200 {
		logger.Error("Got status code %d: %v", statusCode, reply)
		return
	}
	showReply(reply)
}

func showReply(reply Reply) {
	hash := zhash.HashFromMap(reply)
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

func sendCommandWithArgs(args []string) (Reply, int, error) {
	requestUrl := commandUrl
	if len(args) > 1 {
		argString := url.QueryEscape(strings.Join(args[1:], " "))
		requestUrl = commandUrl + argString
	}
	response, err := http.Post(requestUrl, "application/x-www-form-urlencoded",
		nil)
	if err != nil {
		logger.Error(err.Error())
	}
	decoder := json.NewDecoder(response.Body)
	reply := Reply{}
	if err := decoder.Decode(&reply); err != nil {
		return nil, 0, err
	}
	return reply, response.StatusCode, nil
}
