package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Messagess struct {
		InvalidFormat     string `json:"invalidFormat"`
		InvalidDataFormat string `json:"invalidDataFormat"`
		Help              string `json:"help"`
	} `json:"messages"`
	TgBotToken string `json:"tgBotToken"`
}

func LoadConfiguration(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, err
}
