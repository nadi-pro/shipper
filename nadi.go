package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Nadi struct {
		Endpoint   string        `yaml:"endpoint"`
		APIKey     string        `yaml:"apiKey"`
		Token      string        `yaml:"token"`
		Storage    string        `yaml:"storage"`
		Persistent bool          `yaml:"persistent"`
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

func testApiConnectivity(config *Config) error {
	// Create an HTTP client with timeout
	client := &http.Client{
		Timeout: config.Nadi.Timeout,
	}

	// Create an HTTP request
	req, err := http.NewRequest("POST", config.Nadi.Endpoint+"test", nil)
	if err != nil {
		return err
	}

	// Set headers (if required)
	req.Header.Set("Authorization", "Bearer "+config.Nadi.APIKey)
	req.Header.Set("Nadi-Token", config.Nadi.Token)
	req.Header.Set("Accept", "application/vnd.nadi.v1+json")
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print the response
	fmt.Println("Response:", string(body))

	// if connection ok. proceed with reading the files in given log storage
	fmt.Println("HTTP Status Code:", resp.StatusCode)

	return nil
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

	// Call the API endpoint
	err = testApiConnectivity(config)
	if err != nil {
		fmt.Println("Error calling API endpoint:", err)
		return
	}
}
