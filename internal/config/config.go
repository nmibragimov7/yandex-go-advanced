package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Server   *string
	BaseURL  *string
	FilePath *string
	DataBase *string
}

func (c *Config) GetConfig() *Config {
	return c
}

func Init() *Config {
	instance := Config{
		Server:   nil,
		BaseURL:  nil,
		FilePath: nil,
		DataBase: nil,
	}

	flags := flag.NewFlagSet("config", flag.ContinueOnError)

	instance.Server = flags.String("a", ":8080", "Server URL")
	instance.BaseURL = flags.String("b", "http://localhost:8080", "Base URL")
	instance.FilePath = flags.String("f", "./storage.txt", "File path")
	instance.DataBase = flags.String(
		"d",
		"host=localhost user=postgres password=admin dbname=postgres sslmode=disable",
		"Database URL",
	)

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
		instance.FilePath = &envFileStoragePath
	}
	if envDatabase, ok := os.LookupEnv("DATABASE_DSN"); ok {
		instance.DataBase = &envDatabase
	}

	return &instance
}
