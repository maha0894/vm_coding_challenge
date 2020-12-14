package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// Config represents the configuration information.
type Config struct {
	ServiceURL string `json:"service_url"`
	DBPath     string `json:"db_path"`
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
	if Conf.ServiceURL == "" {
		fmt.Println("No ServiceURL has been configured.")
		return errors.New("no ServiceURL has been configured")
	}
	return
}
