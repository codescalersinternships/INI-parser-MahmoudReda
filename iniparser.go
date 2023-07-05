package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config represents the configuration data stored in an INI file
type Config map[string]map[string]string

// ParseINIFile parses the given INI file and returns the configuration data
func ParseINIFile(filePath string) (Config, error) {
	config := make(Config)
	var currentSection string

	// Open the INI file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Found a section header, create a new section in the config
			currentSection = line[1 : len(line)-1]
			config[currentSection] = make(map[string]string)
		} else {
			// Parse key-value pairs within a section
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("malformed line: %s", line)
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[currentSection][key] = value
		}
	}

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}
