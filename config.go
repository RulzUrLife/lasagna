package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Debug    bool `json:"debug"`
	Database struct {
		Name     string `json:"name"`
		User     string `json:"user"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
	} `json:"database"`
}

var (
	Config = getConfig()
)

func getConfig() (config Configuration) {
	configFile, err := os.Open(getConfigPath())
	if err != nil {
		log.Fatalf("Error when reading config file: %s", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Fatalf("Error when unmarshalling the json %s", err)
	}
	return
}

func getConfigPath() (configPath string) {
	if configPath = os.Getenv("LASAGNA_CONFIG"); configPath == "" {
		configPath = "config.json"
	}
	return
}
