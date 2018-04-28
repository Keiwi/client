package main

import (
	"github.com/keiwi/client"
	"github.com/keiwi/utils/log"
	"github.com/keiwi/utils/log/handlers/cli"
	"github.com/keiwi/utils/log/handlers/file"
)

var Version string

func main() {
	fileConfig := file.Config{Folder: "./logs", Filename: "%date%_server.log"}
	log.Log = log.NewLogger(log.DEBUG, []log.Reporter{
		cli.NewCli(),
		file.NewFile(&fileConfig),
	})
	client.Start()
}
