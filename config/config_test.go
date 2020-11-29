package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfigFails1(t *testing.T) {
	t.Run("loads config fails 1", func(t *testing.T) {
		err := LoadConfig("./temp/config.json")
		if err.Error() != "config file error" {
			t.Errorf("got %q, want %q", err.Error(), "config file error")
		}
	})
}

func TestLoadConfigFails2(t *testing.T) {
	t.Run("loads config fails 2", func(t *testing.T) {
		tempConfig := Config{Port: ""}
		d1, _ := json.Marshal(tempConfig)
		err := os.Mkdir("temp", 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("./temp/config.json", d1, 0777)
		if err != nil {
			fmt.Println(err)
		}
		err = LoadConfig("./temp/config.json")
		if !reflect.DeepEqual(tempConfig, Conf) {
			t.Errorf("got %q, want %q", Conf, tempConfig)
		}
		if err.Error() != "no port has been configured" {
			t.Errorf("got %q, want %q", err.Error(), "no port has been configured")
		}
		os.RemoveAll("temp")
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("loads config", func(t *testing.T) {
		tempConfig := Config{Port: "3333"}
		d1, _ := json.Marshal(tempConfig)
		err := os.Mkdir("temp", 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("./temp/config.json", d1, 0777)
		if err != nil {
			fmt.Println(err)
		}
		err = LoadConfig("./temp/config.json")
		if !reflect.DeepEqual(tempConfig, Conf) {
			t.Errorf("got %q, want %q", Conf, tempConfig)
		}
		if err != nil {
			t.Errorf("config loaded with errors %q", err)
		}
		os.RemoveAll("temp")
	})
}
