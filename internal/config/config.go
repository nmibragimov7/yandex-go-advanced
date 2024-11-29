package config

import (
	"flag"
	"os"
)

var (
	Server  *string
	BaseURL *string
)

func Init() string {
	Server = flag.String("a", ":8080", "Server URL")
	BaseURL = flag.String("b", "http://localhost:8080", "Base URL")

	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		Server = &envServerAddress
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		BaseURL = &envBaseURL
	}

	return *Server
}
