package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Server   *string
	BaseURL  *string
	FilePath *string
}

func (c *Config) GetConfig() *Config {
	return c
}

func Init() *Config {
	instance := Config{
		Server:   nil,
		BaseURL:  nil,
		FilePath: nil,
	}

	flags := flag.NewFlagSet("config", flag.ContinueOnError)

	instance.Server = flags.String("a", ":8080", "Server URL")
	instance.BaseURL = flags.String("b", "http://localhost:8080", "Base URL")
	filePath := flags.String("f", "./", "File path")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Printf("ERROR: flag Parse: %s", err.Error())
	}

	if envServerAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		instance.Server = &envServerAddress
	}
	if envBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		instance.BaseURL = &envBaseURL
	}
	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		filePath = &envFileStoragePath
	}

	absPath, err := filepath.Abs(*filePath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	instance.FilePath = &absPath

	return &instance
}
