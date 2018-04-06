package main

import (
	"github.com/keiwi/client"
	"github.com/keiwi/utils/log"
	"github.com/keiwi/utils/log/handlers/cli"
	"github.com/keiwi/utils/log/handlers/file"
)

var Version string

func main() {
	log.Log = log.NewLogger(log.DEBUG, []log.Reporter{
		cli.NewCli(),
		file.NewFile("./logs", "%date%_client.log"),
	})
	client.Start()
}
