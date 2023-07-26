package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"github.com/shirou/gopsutil/disk"
)

type Conf struct {
	Path      string  `yaml:"path"`
	Size      float64 `yaml:"size"`
	SweetSpot float64 `yaml:"sweet_spot"`
}

func main() {
	// Get the path of the current executable
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Cannot determine executable path: %v", err)
	}
	exPath := filepath.Dir(ex)

	// Check if the configuration file exists
	_, err = os.Stat("config.yml")
	if os.IsNotExist(err) {
		// Create the configuration file with default values if it doesn't exist
		conf := &Conf{
			Path:      exPath,
			Size:      100,
			SweetSpot: 10,
		}

		data, err := yaml.Marshal(&conf)
		if err != nil {
			log.Fatalf("Cannot marshal configuration: %v", err)
		}

		err = ioutil.WriteFile("config.yml", data, 0644)
		if err != nil {
			log.Fatalf("Cannot write configuration file: %v", err)
		}

		fmt.Println("A new configuration file has been created with default values. Please modify it according to your needs.")
		os.Exit(1)
	}

	// Read the configuration file
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Cannot read configuration file: %v", err)
	}

	// Parse the configuration file
	var conf Conf
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Cannot unmarshal configuration: %v", err)
	}

	// Get the filesystem information
	diskStat, err := disk.Usage(conf.Path)
	if err != nil {
		log.Fatalf("Failed to get disk usage: %v", err)
	}

	// Calculate the free space in GB
	free := float64(diskStat.Free) / (1024 * 1024 * 1024)
	// Print the available space
	fmt.Printf("Available space: %.4f GB\n", free)

	// Exit with status code 0 if the free space is more or equal than the configured value plus the sweet spot, otherwise exit with status code 1
	if free >= conf.Size+conf.SweetSpot {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
