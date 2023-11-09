package config

import (
	"flag"
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

func ParseConfig() {
	filePath := flag.String("config", "/etc/monomail-sync.yml", "Path of the configuration file in YAML format")
	flag.Parse()

	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Fatalf("Configuration file: %s does not exist, %v\n", *filePath, err)
	}

	viper.SetConfigFile(*filePath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s\n", err)
	}

	err := viper.Unmarshal(&Conf)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v\n", err)
	}
}
