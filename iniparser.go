package iniparser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	ErrInvalidFormat        = errors.New("ini file format isn't valid")
	ErrNoGlobalDataAllowed  = errors.New("Global data isn't supported")
	ErrInvalidFileExtension = errors.New("File extension isn't valid")
	ErrEmpytFile            = errors.New("The file is empty")
	ErrFileNotFound         = errors.New("The system can't find the file")
	ErrSectionNotFound      = errors.New("Section not found")
	ErrKeyNotFound          = errors.New("Key not found")
)

// Config represents the configuration data stored in an INI file
type Config map[string]map[string]string

// LoadFromString loads INI data from a string and returns the configuration data
func LoadFromString(iniData string) (Config, error) {
	config := make(Config)
	var currentSection string

	lines := strings.Split(iniData, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if (line == "" || strings.HasPrefix(line, ";")) && len(lines) == 1 {
			return nil, ErrEmpytFile
		}
		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Found a section header, create a new section in the config
			currentSection = line[1 : len(line)-1]
			//this is allowed?? []
			config[currentSection] = make(map[string]string)
		} else {
			// check if the syntax is not valid
			if strings.HasPrefix(line, "[") || strings.HasSuffix(line, "]") {
				return nil, ErrInvalidFormat
			}

			// Check if global data is present
			if currentSection == "" {
				return nil, ErrNoGlobalDataAllowed
			}

			// Parse key-value pairs within a section
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				return nil, ErrInvalidFormat
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[currentSection][key] = value
		}
	}

	return config, nil
}

// LoadFromFile loads INI data from a file and returns the configuration data
func LoadFromFile(filePath string) (Config, error) {
	if !strings.HasSuffix(filePath, ".ini") {
		return nil, ErrInvalidFileExtension
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, ErrFileNotFound
	}

	return LoadFromString(string(data))
}

// GetSectionNames returns a list of all section names in the configuration
func GetSectionNames(config Config) []string {
	var sectionNames []string
	for section := range config {
		sectionNames = append(sectionNames, section)
	}
	return sectionNames
}

// GetSections returns the configuration data as a dictionary/map of sections and keys
func GetSections(config Config) map[string]map[string]string {
	return config
}

// Get returns the value of a key in a specific section
func Get(config Config, section, key string) (string, error) {
	sectionData, found := config[section]
	if !found {
		return "", ErrSectionNotFound
	}
	value, found := sectionData[key]
	if !found {
		return "", ErrKeyNotFound
	}
	return value, nil
}

// Set sets the value of a key in a specific section
func Set(config Config, section, key, value string) {
	if _, found := config[section]; !found {
		config[section] = make(map[string]string)
	}
	config[section][key] = value
}

// ToString returns the configuration data as a string representation
func ToString(config Config) string {
	var lines []string
	for section, sectionData := range config {
		lines = append(lines, fmt.Sprintf("[%s]", section))
		for key, value := range sectionData {
			lines = append(lines, fmt.Sprintf("%s = %s", key, value))
		}
	}
	return strings.Join(lines, "\n")
}

// SaveToFile saves the configuration data to a file
func SaveToFile(config Config, filePath string) error {
	data := ToString(config)
	return ioutil.WriteFile(filePath, []byte(data), 0644)
}
