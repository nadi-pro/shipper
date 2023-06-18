package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// Config is a Nadi Shipper Configuration based on Yaml
type Config struct {
	Nadi struct {
		Endpoint   string        `yaml:"endpoint"`
		APIKey     string        `yaml:"apiKey"`
		Token      string        `yaml:"token"`
		Storage    string        `yaml:"storage"`
		Persistent bool          `yaml:"persistent"`
		MaxTries   int           `yaml:"maxTries"`
		Timeout    time.Duration `yaml:"timeout"`
		Accept     string        `yaml:"accept"`
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

func callAPIEndpoint(config *Config, endpoint string, payload []byte) error {
	// Create an HTTP client with timeout
	client := &http.Client{
		Timeout: config.Nadi.Timeout,
	}

	// Create an HTTP request
	req, err := http.NewRequest("POST", config.Nadi.Endpoint+endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+config.Nadi.APIKey)
	req.Header.Set("Accept", config.Nadi.Accept)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Nadi-Token", config.Nadi.Token)
	req.Header.Set("Nadi-Transporter-Id", generateTransporterID())

	// Extract transporter ID from payload and add it to request headers
	var payloadData map[string]interface{}
	err = json.Unmarshal(payload, &payloadData)
	if err != nil {
		return err
	}

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

	return nil
}

func generateTransporterID() string {
	randomString := make([]byte, 32)
	_, err := rand.Read(randomString)
	if err != nil {
		panic(err)
	}

	// Compute MD5 hash
	md5Hash := md5.Sum(randomString)
	md5HashString := hex.EncodeToString(md5Hash[:])

	// Compute SHA-1 hash
	sha1Hash := sha1.Sum([]byte(md5HashString))
	sha1HashString := hex.EncodeToString(sha1Hash[:])

	// Return the unique transporter ID
	return sha1HashString
}

func sendJSONFiles(config *Config) {
	// Get the list of JSON files in the directory
	files, err := os.ReadDir(config.Nadi.Storage)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Iterate over the files
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(config.Nadi.Storage, file.Name())

			// Read the content of the JSON file
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}

			// Send the content to the API endpoint
			err = sendToRecordAPI(config, content)
			if err != nil {
				fmt.Println("Error sending to API:", err)
				continue
			}

			// Remove the JSON file
			err = os.Remove(filePath)
			if err != nil {
				fmt.Println("Error removing file:", err)
				continue
			}

			fmt.Println("JSON file processed:", filePath)
		}
	}
}

func sendToRecordAPI(config *Config, payload []byte) error {
	return callAPIEndpoint(config, "record", payload)
}

func verifyAPIEndpoint(config *Config) error {
	return callAPIEndpoint(config, "verify", nil)
}

func main() {
	// Parse command-line arguments
	configPath := flag.String("config", "nadi.yaml", "path to config file")
	verifyFlag := flag.Bool("verify", false, "verify API endpoint")
	testFlag := flag.Bool("test", false, "test API endpoint")
	recordFlag := flag.Bool("record", false, "test API endpoint")
	flag.Parse()

	fmt.Println("Nadi Ship set sailing at " + time.Now().Format("2006-01-02 15:04:05"))

	// Load the configuration from YAML
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	// Test the API endpoint if -test flag is provided
	if *testFlag {
		err = callAPIEndpoint(config, "test", nil)
		if err != nil {
			fmt.Println("Error calling API endpoint:", err)
			return
		}
		return
	}

	// Verify the API endpoint if -verify flag is provided
	if *verifyFlag {
		err = verifyAPIEndpoint(config)
		if err != nil {
			fmt.Println("Error verifying API endpoint:", err)
			return
		}
		return
	}

	// Check for JSON files in the storage directory and send them to the record API
	if *recordFlag {
		sendJSONFiles(config)
		return
	}

	fmt.Println(generateTransporterID())

	fmt.Println("Nadi Ship successfully deliver the goods at " + time.Now().Format("2006-01-02 15:04:05"))
}
