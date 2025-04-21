package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// Config - config struct
type Config struct {
	Server   *string `json:"server"`
	BaseURL  *string `json:"base_url"`
	FilePath *string `json:"file_path"`
	DataBase *string `json:"data_base"`
	HTTPS    *bool   `json:"https"`
	Config   *string
}

// parseJSON - parse json instance
func parseJSON(config *Config) error {
	if config.Config == nil {
		return nil
	}

	file, err := os.Open(*config.Config)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("failed to close config file: %s", err.Error())
		}
	}(file)

	var jsonConf Config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&jsonConf); err != nil {
		return fmt.Errorf("failed to decode json file: %w", err)
	}
	if jsonConf.Server != nil && *jsonConf.Server != "" && *config.Server == "" {
		config.Server = jsonConf.Server
	}
	if jsonConf.BaseURL != nil && *jsonConf.BaseURL != "" && *config.BaseURL == "" {
		config.BaseURL = jsonConf.BaseURL
	}
	if jsonConf.FilePath != nil && *jsonConf.FilePath != "" && *config.FilePath == "" {
		config.FilePath = jsonConf.FilePath
	}
	if jsonConf.DataBase != nil && *jsonConf.DataBase != "" && *config.DataBase == "" {
		config.DataBase = jsonConf.DataBase
	}
	if jsonConf.HTTPS != nil && *jsonConf.HTTPS && !*config.HTTPS {
		config.HTTPS = jsonConf.HTTPS
	}

	return nil
}

// Init - initialize config instance
func Init() *Config {
	instance := Config{
		Server:   nil,
		BaseURL:  nil,
		FilePath: nil,
		DataBase: nil,
		HTTPS:    nil,
		Config:   nil,
	}

	flags := flag.NewFlagSet("config", flag.ContinueOnError)

	instance.Server = flags.String("a", ":8080", "Server URL")
	instance.BaseURL = flags.String("b", "http://localhost:8080", "Base URL")
	instance.FilePath = flags.String("f", "", "File path") // ./storage.txt
	instance.DataBase = flags.String(
		"d",
		"",
		"Database URL",
	) // host=localhost user=postgres password=admin dbname=postgres sslmode=disable
	instance.HTTPS = flags.Bool("s", false, "Enable HTTPS")
	instance.Config = flags.String("c", "", "Config path")

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
	if envHTTPS, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		if envHTTPS == "true" {
			value := true
			instance.HTTPS = &value
		}
	}
	if envConfigPath, ok := os.LookupEnv("CONFIG"); ok {
		instance.Config = &envConfigPath
	}

	if err = parseJSON(&instance); err != nil {
		log.Printf("failed to parse json: %s", err.Error())
	}

	fmt.Println("Server", *instance.Server)
	fmt.Println("BaseURL", *instance.BaseURL)
	fmt.Println("FilePath", *instance.FilePath)
	fmt.Println("DataBase", *instance.DataBase)
	fmt.Println("HTTPS", *instance.HTTPS)
	fmt.Println("Config", *instance.Config)

	return &instance
}
