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

// parseFlags - parse flags instance
func parseFlags() *Config {
	var instance Config

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

	return &instance
}

// parseEnv - parse env instance
func parseEnv() *Config {
	var instance Config

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

	return &instance
}

// parseJSON - parse json instance
func parseJSON(config *Config) (error, *Config) {
	if config.Config == nil {
		return nil, config
	}

	file, err := os.Open(*config.Config)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err), nil
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
		return err, nil
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

	return nil, config
}

// Init - initialize config instance
func Init() *Config {
	var instance Config

	instance = *parseFlags()
	if d := *parseEnv(); d.Server != nil {
		instance = *parseEnv()
	}
	if err, cnf := parseJSON(&instance); err == nil {
		instance = *cnf
	}

	return &instance
}
