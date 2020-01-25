package tools

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Database struct {
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"database"`

	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
}

func LoadConfiguration(filename string) (*Config, error) {
	config := &Config{}

	configFile, err := os.Open(filename)
	defer func() {
		if err := configFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
