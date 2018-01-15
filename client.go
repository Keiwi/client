package client

import (
	"io"
	"net"
	"strings"

	"encoding/json"

	"github.com/apex/log"
	"github.com/keiwi/utils"
)

var (
	version string
	updater UpdateService
)

// ErrorResponse is the struct when sending a simple error back
type ErrorResponse struct {
	Error string `json:"error"`
}

// StartTCPServer starts a TCP server and waits for requests
func StartTCPServer(connIP, connPort, connType string) {
	tcp, err := net.Listen(connType, connIP+":"+connPort)
	if err != nil {
		utils.Log.WithField("error", err).Fatal("Error listening on TCP")
	}

	// Close the tcp request when done
	defer tcp.Close()

	utils.Log.WithFields(log.Fields{
		"IP":   connIP,
		"Port": connPort,
	}).Info("Started listening for TCP requests")
	for {
		// Wait for incoming requests
		conn, err := tcp.Accept()
		if err != nil {
			utils.Log.WithField("error", err).Error("Error accepting TCP request")
			return
		}

		// Take care of the request in a gorotuine
		go handleRequest(conn)
	}
}

// Start starts the whole client, takes care of the config,
// starts the TCP server etc.
func Start(v string) {
	version = v
	utils.Log.WithField("Version", version).Info("Starting MSTT-Monitor client")

	updater = UpdateService{
		Version:    version,
		Identifier: "mstt-client-windows-",
	}

	StartTCPServer("0.0.0.0", "3333", "tcp")
}

// handleRequest handles the incoming TCP requests,
// parses the command and responds to the request
func handleRequest(conn io.ReadWriteCloser) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err.Error() == "EOF" {
			return
		}
		utils.Log.WithField("error", err).Error("Error reading TCP request")
		return
	}
	cmd := ParseCommand(string(buf[:n]))

	var resp interface{}

	switch strings.ToLower(cmd.Name) {
	case "check_memory":
		resp = MemoryCheck(cmd)
	case "check_disc":
		resp = DiscCheck(cmd)
	case "check_cpu":
		resp = CPUCheck(cmd)
	case "uptime":
		resp = UptimeCheck(cmd)
	case "info":
		resp = InfoCheck(cmd)
	case "file":
		resp = FileCheck(cmd)
	case "update":
		resp = UpdateCheck(cmd)
	case "netusage":
		resp = NetworkCheck(cmd)
	default:
		resp = ErrorResponse{Error: "Unknown command"}
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		utils.Log.WithField("error", err).Error("Error parsing respBody")
		return
	}

	conn.Write(respBody)
	conn.Close()
}
