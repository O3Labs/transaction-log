package config

import (
	"encoding/json"
	"os"
)

func LoadConfigurationFile(file string) (Configuration, error) {
	configuration := Configuration{}
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return configuration, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configuration)
	return configuration, nil
}

type Configuration struct {
	SeedList    []string `json:"seedList"`
	RPCSeedList []string `json:"rpcSeedList"`
	Magic       int      `json:"magic"` //network ID.
}
