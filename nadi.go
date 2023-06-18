package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Nadi struct {
		Endpoint   string `yaml:"endpoint"`
		APIKey     string `yaml:"apiKey"`
		Token      string `yaml:"token"`
		Storage    string `yaml:"storage"`
		Persistent bool   `yaml:"persistent"`
		MaxTries   int           `yaml:"maxTries"`
		Timeout    time.Duration `yaml:"timeout"`
	} `yaml:"nadi"`
}

func loadConfig(filename string) (*Config, error) {
	// Read the YAML file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse the YAML data into the Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	// Parse command-line arguments
	configPath := flag.String("config", "nadi.yaml", "Path to the configuration file")
	flag.Parse()

	// Load the configuration from the YAML file
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	// Access the configuration values
	fmt.Println("Endpoint:", config.Nadi.Endpoint)
	fmt.Println("API Key:", config.Nadi.APIKey)
	fmt.Println("Token:", config.Nadi.Token)
	fmt.Println("Storage:", config.Nadi.Storage)
	fmt.Println("Persistent:", config.Nadi.Persistent)
	fmt.Println("Max Tries:", config.Nadi.MaxTries)
	fmt.Println("Timeout:", config.Nadi.Timeout)
}
