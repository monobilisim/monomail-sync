package config

import (
	"imap-sync/logger"
	"os"

	"github.com/spf13/viper"
)

var log = logger.Log

type Config struct {
	Language string
	Port     string

	DatabaseInfo struct {
		AdminName    string
		AdminPass    string
		DatabasePath string
	}

	SourceAndDestination struct {
		SourceServer      string
		SourceMail        string
		DestinationServer string
		DestinationMail   string
	}

	Email struct {
		SMTPHost     string
		SMTPPort     string
		SMTPFrom     string
		SMTPUser     string
		SMTPPassword string
	}
}

var Conf Config
var configFilePath string

func SetConfigPath(path string) {
	configFilePath = path
}

func ParseConfig() {
	if configFilePath == "" {
		configFilePath = "/etc/monomail-sync.yml"
	}

	setDefaults()

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Warnf("Config file not found, using defaults: %v", err)
		return
	}

	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Error reading config file, using defaults: %v", err)
		return
	}

	viper.Unmarshal(&Conf)
}

func setDefaults() {
	Conf.Language = "en"
	Conf.Port = "8000"
	Conf.DatabaseInfo.AdminName = "admin"
	Conf.DatabaseInfo.AdminPass = "admin"
	Conf.DatabaseInfo.DatabasePath = "./db.db"
	Conf.SourceAndDestination.SourceServer = "imap.example.com"
	Conf.SourceAndDestination.SourceMail = "@example.com"
	Conf.SourceAndDestination.DestinationServer = "imap.example.com"
	Conf.SourceAndDestination.DestinationMail = "@example.com"
}
