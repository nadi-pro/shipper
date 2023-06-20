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

type Config struct {
	Nadi struct {
		Endpoint    string        `yaml:"endpoint"`
		APIKey      string        `yaml:"apiKey"`
		Token       string        `yaml:"token"`
		Storage     string        `yaml:"storage"`
		Persistent  bool          `yaml:"persistent"`
		MaxTries    int           `yaml:"maxTries"`
		Timeout     time.Duration `yaml:"timeout"`
		Accept      string        `yaml:"accept"`
		TrackerFile string        `yaml:"trackerFile"`
	} `yaml:"nadi"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type FileStatus int

const (
	StatusPending FileStatus = iota
	StatusSent
	StatusFailed
)

type FileTracker struct {
	Status FileStatus
	Tries  int
}

type TrackerMap map[string]FileTracker

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

	// Set Payload
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

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		// Try to parse the response as an error message
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err == nil && errorResponse.Message != "" {
			return fmt.Errorf("API request failed with status code: %d, Response: %s", resp.StatusCode, errorResponse.Message)
		}

		// Fallback to returning the raw response body
		return fmt.Errorf("API request failed with status code: %d, Response: %s", resp.StatusCode, string(body))
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
	files, err := ioutil.ReadDir(config.Nadi.Storage)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Create a tracker map to store the status of each file
	trackerMap := make(TrackerMap)

	// Load the tracker data from a file (including creating the file if it doesn't exist)
	trackerFilePath := filepath.Join(config.Nadi.Storage, config.Nadi.TrackerFile)
	loadTrackerData(trackerFilePath, &trackerMap)

	// Iterate over the files
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(config.Nadi.Storage, file.Name())

			// Check if the file is already sent or failed
			if trackerMap[file.Name()].Status != StatusPending {
				continue
			}

			// Read the JSON file content
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}

			// Call the API endpoint with the JSON content
			err = callAPIEndpoint(config, "record", content)
			if err != nil {
				fmt.Println("Error calling API:", err)
				// Increment the number of tries for the file
				tracker := trackerMap[file.Name()]
				tracker.Tries++

				// Mark the file as failed if max tries exceeded
				if tracker.Tries > config.Nadi.MaxTries {
					tracker.Status = StatusFailed
				}
				trackerMap[file.Name()] = tracker
				continue
			}

			// Mark the file as sent
			tracker := trackerMap[file.Name()]
			tracker.Status = StatusSent
			trackerMap[file.Name()] = tracker

			// Remove the JSON file if not persistent
			if !config.Nadi.Persistent {
				err := os.Remove(filePath)
				if err != nil {
					fmt.Println("Error removing file:", err)
				}
			}
		}
	}

	// Save the tracker data to a file
	saveTrackerData(trackerFilePath, trackerMap)
}

func verifyAPIEndpoint(config *Config) {
	err := callAPIEndpoint(config, "verify", nil)
	if err != nil {
		fmt.Println("API verification failed:", err)
		return
	}

	fmt.Println("API verification successful.")
}

func testAPIEndpoint(config *Config) {
	err := callAPIEndpoint(config, "test", nil)
	if err != nil {
		fmt.Println("Error calling test API:", err)
		return
	}

	fmt.Println("Test API call successful.")
}

func loadTrackerData(filepath string, trackerMap *TrackerMap) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		// Create an empty tracker map if the file doesn't exist
		if os.IsNotExist(err) {
			*trackerMap = make(TrackerMap)
			return
		}

		fmt.Println("Error loading tracker data:", err)
		return
	}

	err = json.Unmarshal(data, trackerMap)
	if err != nil {
		fmt.Println("Error loading tracker data:", err)
		*trackerMap = make(TrackerMap)
	}
}

func saveTrackerData(filePath string, trackerMap TrackerMap) {
	// Convert the tracker map to JSON
	data, err := json.Marshal(trackerMap)
	if err != nil {
		fmt.Println("Error encoding tracker data:", err)
		return
	}

	// Write the tracker data to a file
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Println("Error writing tracker file:", err)
	}
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
		testAPIEndpoint(config)
		return
	}

	// Verify the API endpoint if -verify flag is provided
	if *verifyFlag {
		verifyAPIEndpoint(config)
		return
	}

	// Check for JSON files in the storage directory and send them to the record API
	if *recordFlag {
		sendJSONFiles(config)
		return
	}

	fmt.Println("Nadi Ship successfully deliver the goods at " + time.Now().Format("2006-01-02 15:04:05"))
}
