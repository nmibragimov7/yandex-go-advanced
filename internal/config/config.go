package config

import (
	"flag"
	"os"
)

type Config struct {
	server  *string
	baseURL *string
}

var globalConfig = Config{
	server:  nil,
	baseURL: nil,
}

func GetBaseURL() *string {
	return globalConfig.baseURL
}

func Init() string {
	globalConfig.server = flag.String("a", ":8080", "Server URL")
	globalConfig.baseURL = flag.String("b", "http://localhost:8080", "Base URL")

	flag.Parse()

	if envServerAddress, exists := os.LookupEnv("SERVER_ADDRESS"); exists {
		globalConfig.server = &envServerAddress
	}
	if envBaseURL, exists := os.LookupEnv("BASE_URL"); exists {
		globalConfig.baseURL = &envBaseURL
	}

	return *globalConfig.server
}
