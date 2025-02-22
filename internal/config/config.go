package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Server     *string
	BaseURL    *string
	FilePath   *string
	DataBase   *string
	SercretKey *string
}

func Init() *Config {
	instance := Config{
		Server:     nil,
		BaseURL:    nil,
		FilePath:   nil,
		DataBase:   nil,
		SercretKey: nil,
	}

	flags := flag.NewFlagSet("config", flag.ContinueOnError)

	instance.Server = flags.String("a", ":8080", "Server URL")
	instance.BaseURL = flags.String("b", "http://localhost:8080", "Base URL")
	instance.FilePath = flags.String("f", "./storage.txt", "File path") // ./storage.txt
	instance.DataBase = flags.String(
		"d",
		"",
		"Database URL",
	) // host=localhost user=postgres password=admin dbname=postgres sslmode=disable
	instance.SercretKey = flags.String("s", "secret_key", "Cookie secret key")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Printf("failed to parse flags: %s", err.Error())
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
	if envDatabase, ok := os.LookupEnv("SECRET_KEY"); ok {
		instance.SercretKey = &envDatabase
	}

	return &instance
}
