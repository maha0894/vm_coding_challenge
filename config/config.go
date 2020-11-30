package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// Config represents the configuration information.
type Config struct {
	Port   string `json:"port"`
	DBPath string `json:"db_path"`
}

// Conf contains the initialized configuration struct
var Conf Config

// LoadConfig loads the configuration from the specified filepath
func LoadConfig(filepath string) (err error) {
	// Get the config file
	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("File error: ", err)
		return errors.New("config file error")
	}
	json.Unmarshal(configFile, &Conf)
	if Conf.Port == "" {
		fmt.Println("No port has been configured.")
		return errors.New("no port has been configured")
	}
	return
}
