package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type msHost struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

var (
	MsHost msHost
)

func init() {
	iniPath := "config/config.ini"
	if args := os.Args; len(args) > 1 {
		iniPath = args[1]
	}

	iniFile, err := ini.Load(iniPath)
	if err != nil {
		log.Fatalf("load %s error:%s\n", iniPath, err.Error())
	}

	//mysql
	database := iniFile.Section("mysql")
	MsHost.Host = database.Key("Host").String()
	MsHost.Port = database.Key("Port").String()
	MsHost.User = database.Key("User").String()
	MsHost.Pass = database.Key("Pass").String()
	MsHost.Name = database.Key("Name").String()
}
