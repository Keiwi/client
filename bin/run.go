package main

import (
	"github.com/Keiwi/client"
	"github.com/Keiwi/utils"
	"github.com/apex/log"
)

var Version string

func main() {
	utils.Log = utils.NewLogger(log.DebugLevel, &utils.LoggerConfig{
		Dirname: "./logs",
		Logname: "%date%_client.log",
	})
	client.Start(Version)
}
