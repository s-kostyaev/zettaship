package main

import (
	"github.com/BurntSushi/toml"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
)

const (
	configPath = "/etc/zettaship.toml"
)

var (
	logfile   = os.Stderr
	formatter = logging.MustStringFormatter(
		"%{time:15:04:05.000000} %{pid} %{level:.8s} %{message}")
	loglevel = logging.INFO
	logger   = logging.MustGetLogger("fake-zfs")
	config   = getConfig(configPath)
)

type Config struct {
	ServerUrl string
}

func setupLogger() {
	logging.SetBackend(logging.NewLogBackend(logfile, "", 0))
	logging.SetFormatter(formatter)
	logging.SetLevel(loglevel, logger.Module)
}

func getConfig(configPath string) *Config {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	config := Config{}
	_, err = toml.Decode(string(buf), &config)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return &config
}
