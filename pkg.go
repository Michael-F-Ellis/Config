package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// In this package, a configuration is represented as a map[string]any and serialized to a file in JSON format.
// Package config provides methods to read, write and update configurations.

type Config map[string]any

// Write writes a configuration to a file.
func Write(filepath string, config Config) (err error) {
	text, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("Unable to serialize config: %v", err)
	}
	err = os.WriteFile(filepath, text, 0644)
	if err != nil {
		return fmt.Errorf("Unable to write config file: %v", err)
	}
	return
}

// Read reads a configuration from a file and rewrites the target configuration.
func Read(filepath string, target Config) (err error) {
	text, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("Unable to read config file: %v", err)
	}
	err = json.Unmarshal(text, &target)
	if err != nil {
		return fmt.Errorf("Unable to parse config file: %v", err)
	}
	return
}
