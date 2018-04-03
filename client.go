package client

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/keiwi/client/commands"
	"github.com/keiwi/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	configType     string
	conn           net.Conn
	commandHandler *commands.CommandHandler
)

// Start will start the whole client process
func Start() {
	ReadConfig()

	utils.Log.Info("Initializing all of the commands")
	commandHandler = commands.NewCommandHandler()

	utils.Log.Info("Starting keiwi Monitor Client")

	utils.Log.Info("Starting discovery server")
	go StartDiscovery()

	utils.Log.Info("Starting to connect to server")
	Connect()
}

// ReadConfig will try to find the config and read, if config file
// does not exists it will create one with default options
func ReadConfig() {
	configType = os.Getenv("KeiwiConfigType")
	if configType == "" {
		configType = "json"
	}
	viper.SetConfigType(configType)

	viper.SetConfigFile("config." + configType)
	viper.AddConfigPath(".")

	viper.SetDefault("log_dir", "./logs")
	viper.SetDefault("log_syntax", "%date%_server.log")
	viper.SetDefault("log_level", "info")

	viper.SetDefault("server_ip", "")
	viper.SetDefault("password", "")
	viper.SetDefault("interval", 600)
	viper.SetDefault("certificate_path", "./server.crt")

	if err := viper.ReadInConfig(); err != nil {
		utils.Log.Debug("Config file not found, saving default")
		if err = viper.WriteConfigAs("config." + configType); err != nil {
			utils.Log.WithField("error", err.Error()).Fatal("Can't save default config")
		}
	}

	level := strings.ToLower(viper.GetString("log_level"))
	utils.Log = utils.NewLogger(utils.NameToLevel[level], &utils.LoggerConfig{
		Dirname: viper.GetString("log_dir"),
		Logname: viper.GetString("log_syntax"),
	})
}

func Connect() {
	caCert, err := ioutil.ReadFile(viper.GetString("certificate_path"))
	if err != nil {
		utils.Log.WithField("error", err.Error()).Fatal("Can't read pem file")
		return
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	conf := &tls.Config{
		RootCAs: certPool,
	}

	for {
		con, err := tls.Dial("tcp", viper.GetString("server_ip"), conf)
		if err != nil {
			utils.Log.WithField("ip", viper.GetString("server_ip")).WithError(err).Error("can't connect to server")

			utils.Log.Infof("failed to connect to server, trying again in %d seconds", viper.GetInt("interval"))
			time.Sleep(time.Second * time.Duration(viper.GetInt("interval")))
			continue
		}
		utils.Log.WithField("IP", con.RemoteAddr().String()).Info("connected to the server")
		conn = con

		utils.Log.Info("Initializing the handshake")
		err = Handshake(conn)
		if err != nil {
			utils.Log.WithError(err).Error("handshake failed")

			utils.Log.Infof("failed to connect to server, trying again in %d seconds", viper.GetInt("interval"))
			time.Sleep(time.Second * time.Duration(viper.GetInt("interval")))
			continue
		}
		utils.Log.Info("Handshake successful")

		for {
			r := bufio.NewReader(conn)
			msg, err := r.ReadString('\n')
			if err != nil {
				utils.Log.WithField("error", err).Error("error reading TLS message")
				break
			}

			cmd := commands.ParseCommand(msg)
			out := commandHandler.RunCommand(cmd)

			data := ""
			if out.Error() != "" {
				data = `{"error": "` + out.Error() + `"}`
			} else {
				b, err := json.Marshal(out.Message())
				if err != nil {
					data = `{"error": "` + err.Error() + `"}`
				} else {
					data = `{"message": ` + string(b) + `}`
				}
			}

			_, err = fmt.Fprintln(conn, data)
			if err != nil {
				utils.Log.WithError(err).Error("error writing TLS message")
				continue
			}
		}
	}
}

func Handshake(conn net.Conn) error {
	_, err := fmt.Fprintln(conn, viper.GetString("password"))
	if err != nil {
		return err
	}

	r := bufio.NewReader(conn)
	accepted, err := r.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "connection disconnected")
	}
	accepted = strings.TrimSpace(accepted)

	if accepted != "accepted" {
		return errors.New("invalid password")
	}
	return nil
}
