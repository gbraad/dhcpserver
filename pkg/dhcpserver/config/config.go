package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var Config *ConfigType

// Create new object with data if file exists or
// Create json file and return object if doesn't exists
func NewConfig(path string) (*ConfigType, error) {
	cfg := &ConfigType{StaticAssignments: []StaticAssignmentsConfigType{}, ServerConfiguration: ServerConfigType{}}
	cfg.ServerConfiguration.Options = make(map[string]string)
	cfg.FilePath = path

	// Check json file existence
	_, err := os.Stat(cfg.FilePath)
	if os.IsNotExist(err) {
		if errWrite := cfg.Write(); errWrite != nil {
			return nil, errWrite
		}
	} else {
		if errRead := cfg.read(); errRead != nil {
			return nil, errRead
		}
	}

	return cfg, nil
}

func (cfg *ConfigType) Write() error {
	jsonData, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(cfg.FilePath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

func (cfg *ConfigType) Delete() error {
	if err := os.Remove(cfg.FilePath); err != nil {
		return err
	}

	return nil
}

func (cfg *ConfigType) read() error {
	raw, err := ioutil.ReadFile(cfg.FilePath)
	if err != nil {
		return err
	}

	if !json.Valid(raw) {
		return errors.New("Invalid JSON")
	}

	json.Unmarshal(raw, &cfg)

	return nil
}
