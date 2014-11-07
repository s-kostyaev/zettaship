package main

func main() {
	setupLogger()
	createRequest()
	addArgs()
	addCookie()
	parseMessage(sendRequest())
}
