package client

import (
	"io/ioutil"
	"net"

	"github.com/apex/log"
	"github.com/keiwi/utils"
	"github.com/spf13/viper"
)

// StartDiscovery will spawn a small TCP server that waits for any incoming discovery messages
func StartDiscovery() {
	tcp, err := net.Listen("tcp", "0.0.0.0:3333")
	if err != nil {
		utils.Log.WithField("error", err).Fatal("error listening on TCP")
	}

	// Close the tcp request when done
	defer tcp.Close()

	utils.Log.WithFields(log.Fields{
		"IP":   "0.0.0.0",
		"Port": "3333",
	}).Info("started listening for auto discovery requests")

	for {
		// Wait for incoming requests
		conn, err := tcp.Accept()
		if err != nil {
			utils.Log.WithField("error", err).Error("error accepting TCP request")
			return
		}

		// Take care of the request in a gorotuine
		go handleDiscovery(conn)
	}
}

// handleDiscovery will handle all TCP requests related to discovery
func handleDiscovery(conn net.Conn) {
	defer conn.Close()
	msg, err := ioutil.ReadAll(conn)
	if err != nil {
		utils.Log.WithField("error", err).Error("error reading TCP request")
		return
	}

	if string(msg) != "discovery" {
		utils.Log.Error("invalid discovery message")
		return
	}

	ip := conn.LocalAddr().String()
	viper.Set("server_ip", ip)
	if err = viper.WriteConfigAs("config." + configType); err != nil {
		utils.Log.WithField("error", err.Error()).Fatal("Can't save config")
		return
	}

	// RestartConnection
}
