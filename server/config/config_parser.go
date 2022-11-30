package config

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// The data struct to decode config.json
type Config struct {
	Master     bool   `json:"master"`
	TLS        bool   `json:"tls"`
	Port       string `json:"port"`
	APIKey     string `json:"api_key"`
	APIKeyHash [32]byte
	Fullchain  string `json:"fullchain"`
	PrivKey    string `json:"priv_key"`
	Clients    []struct {
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
	var content []byte

	// get path of running executable
	filename, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(filename)
	// build file path as wd + '/config.json'
	filePath := exPath + "/config.json"
	content, err = ioutil.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(content, &result)
	result.APIKeyHash = sha256.Sum256([]byte(result.APIKey))
	if err != nil {
		return Config{}, err
	}

	return result, nil
}
