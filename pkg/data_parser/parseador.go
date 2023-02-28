package parseador

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Connection       string   `json:"connection"`
	Enabled_Backends []string `json:"enabled_backends"`
	NumberOfBackends string   `json:"numberofbackends"`
}

func Parseador(file string) {
	// Read the contents of the configuration file into a byte slice
	data, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		logrus.Fatal("Can't read file", err)
	}

	// Split the byte slice into a slice of strings, with each string representing a line in the file
	lines := strings.Split(string(data), "\n")
	logrus.Info("Lines split in slice of strings")
	logrus.Info("Number of lines: ", len(lines))

	// Initialize an empty Config struct
	var config Config

	// Iterate through the lines in the file
	for _, line := range lines {
		//logrus.Info("Untreated lines: ", lines)

		// Skip empty lines
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		// Split the line into a key-value pair
		parts := strings.Split(line, "=")
		logrus.Info(line, " --- ", "output: ", parts)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set the appropriate field in the Config struct based on the key
		var x int
		switch key {
		case "enabled_backends":
			config.Enabled_Backends = append(config.Enabled_Backends, value)
			config.NumberOfBackends = fmt.Sprint(x + 1)
		}
	}

	// Marshal the Config struct into a JSON byte slice
	configJSON, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(configJSON))
}
