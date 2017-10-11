/*
Package dbshield implements the database firewall functionality
*/
package dbshield

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/qiwihui/DBShield/dbshield/config"
	"github.com/qiwihui/DBShield/dbshield/httpserver"
	"github.com/qiwihui/DBShield/dbshield/logger"
)

//Version of the library
var Version = "1.0.0-beta4-qiwihui"

var configFile string

//SetConfigFile of DBShield
func SetConfigFile(cf string) error {
	configFile = cf
	err := config.ParseConfig(configFile)
	if err != nil {
		return err
	}
	return postConfig()
}

//ShowConfig writes parsed config file as JSON to STDUT
func ShowConfig() error {
	confJSON, err := json.MarshalIndent(config.Config, "", "    ")
	fmt.Println(string(confJSON))
	return err
}

//Purge local database
func Purge() error {
	return config.Config.LocalDB.Purge()
}

//Patterns lists the captured patterns
func Patterns() (count int) {
	return config.Config.LocalDB.Patterns()
}

//Abnormals detected querties
func Abnormals() (count int) {
	return config.Config.LocalDB.Abnormals()
}

//RemovePattern deletes a pattern from captured patterns DB
func RemovePattern(pattern string) error {
	return config.Config.LocalDB.DeletePattern([]byte(pattern))
}

func postConfig() (err error) {

	config.Config.DB, err = dbNameToStruct(config.Config.DBType)
	if err != nil {
		return err
	}

	tmpDBMS, _ := generateDBMS()
	if config.Config.ListenPort == 0 {
		config.Config.ListenPort = tmpDBMS.DefaultPort()
	}
	if config.Config.TargetPort == 0 {
		config.Config.TargetPort = tmpDBMS.DefaultPort()
	}
	return
}

func mainListner() error {
	if config.Config.HTTP {
		proto := "http"
		if config.Config.HTTPSSL {
			proto = "https"
		}
		logger.Infof("Web interface on %s://%s/", proto, config.Config.HTTPAddr)
		go httpserver.Serve()
	}
	serverAddr, _ := net.ResolveTCPAddr("tcp", config.Config.TargetIP+":"+strconv.Itoa(int(config.Config.TargetPort)))
	l, err := net.Listen("tcp", config.Config.ListenIP+":"+strconv.Itoa(int(config.Config.ListenPort)))
	if err != nil {
		return err
	}
	// Close the listener when the application closes.
	defer l.Close()

	for {
		// Listen for an incoming connection.
		listenConn, err := l.Accept()
		if err != nil {
			logger.Warningf("Error accepting connection: %v", err)
			continue
		}
		go handleClient(listenConn, serverAddr)
	}
}

//Start the proxy
func Start() (err error) {

	initLogging()
	logger.Infof("Config file: %s", configFile)
	logger.Infof("Listening: %s:%v",
		config.Config.ListenIP,
		config.Config.ListenPort)
	logger.Infof("Backend: %s (%s:%v)",
		config.Config.DBType,
		config.Config.TargetIP,
		config.Config.TargetPort)
	logger.Infof("Protect: %v", !config.Config.Learning)
	logger.Infof("Recording queries: %v", config.Config.LocalQueryRecord)
	go mainListner()
	signalHandler()
	return nil
}
