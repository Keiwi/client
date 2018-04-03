package commands

import (
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
)

//region UptimeOutput

// NetworkInterface contains information about a network interface
type NetworkInterface struct {
	Name string   `json:"name"`
	IPs  []string `json:"ips"` // Can be both ipv4 and ipv6
}

// InfoOutput contains information about uptime of a client
type InfoOutput struct {
	Hostname      string
	OS            string
	Platform      string
	ClientVersion string
	Interfaces    []NetworkInterface
}

func (InfoOutput) Error() string { return "" }

func (i InfoOutput) Message() OutputMessage {
	return map[string]interface{}{
		"hostname":       i.Hostname,
		"os":             i.OS,
		"platform":       i.Platform,
		"client_version": i.ClientVersion,
		"interfaces":     i.Interfaces,
	}
}

//endregion

//region InfoCommand

// InfoCommand gathers information about the client
type InfoCommand struct {
	Version string
}

func (i InfoCommand) Run(cmd Command) Output {
	info, err := host.Info()
	if err != nil {
		return ErrorOutput{err: err.Error()}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return ErrorOutput{err: err.Error()}
	}

	resp := InfoOutput{
		Hostname:      info.Hostname,
		OS:            info.OS,
		Platform:      info.Platform,
		ClientVersion: i.Version,
	}

	for _, inter := range interfaces {
		for _, up := range inter.Flags {
			if up != "up" {
				continue
			}

			networkInterface := NetworkInterface{Name: inter.Name}
			for _, addr := range inter.Addrs {
				networkInterface.IPs = append(networkInterface.IPs, addr.Addr)
			}
			resp.Interfaces = append(resp.Interfaces, networkInterface)
		}
	}

	return resp
}

func (InfoCommand) Name() string { return "info" }
func (InfoCommand) Description() string {
	return "returns info about the client"
}
func (InfoCommand) Usage() string { return `` }

//endregion
