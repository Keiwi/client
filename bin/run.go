package main

import (
	"github.com/apex/log"
	"github.com/keiwi/client"
	"github.com/keiwi/utils"
)

var Version string

func main() {
	utils.Log = utils.NewLogger(log.DebugLevel, &utils.LoggerConfig{
		Dirname: "./logs",
		Logname: "%date%_client.log",
	})
	client.Start(Version)
}
