package config

import (
	"encoding/json"
	"io/ioutil"
	"runtime"
)

// The data struct to decode config.json
type Config struct {
	Master    bool   `json:"master"`
	TLS       bool   `json:"tls"`
	Port      string `json:"port"`
	APIKey    string `json:"api_key"`
	Fullchain string `json:"fullchain"`
	PrivKey   string `json:"priv_key"`
	Clients   []struct {
		Alias string `json:"alias"`
		IP    string `json:"ip"`
		MAC   string `json:"mac"`
		Os    []struct {
			Name   string `json:"name"`
			BootID string `json:"boot_id"`
		} `json:"os,omitempty"`
	} `json:"clients:"`
}

func ParseConfig() (Config, error) {

	var result Config

	// if OS is windows, use windows path
	// if OS is linux, use linux path
	var content []byte
	var err error
	os := runtime.GOOS
	if os == "windows" {
		content, err = ioutil.ReadFile("C:\\wolservice\\config\\config.json")
	} else if os == "linux" {
		content, err = ioutil.ReadFile("./config/config.json")
	}

	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		return Config{}, err
	}

	return result, nil
}
