package config

import (
	"flag"
	"os"
)

type Config struct {
	Server  *string
	BaseURL *string
}

var globalConfig = Config{
	Server:  nil,
	BaseURL: nil,
}

func GetConfig() Config {
	return globalConfig
}

func Init() string {
	globalConfig.Server = flag.String("a", ":8080", "Server URL")
	globalConfig.BaseURL = flag.String("b", "http://localhost:8080", "Base URL")

	flag.Parse()

	if envServerAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		globalConfig.Server = &envServerAddress
	}
	if envBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		globalConfig.BaseURL = &envBaseURL
	}

	return *globalConfig.Server
}
