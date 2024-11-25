package config

import "flag"

var (
	Server  *string
	BaseURL *string
)

func Init() {
	Server = flag.String("a", ":8080", "Server URL")
	BaseURL = flag.String("b", "http://localhost:8080", "Base URL")
}
