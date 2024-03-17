// Package config provides methods to read, write and update configurations.
// The intended use case is in programs that need to make use of API's that
// accept POST requests with JSON payloads.  The package is designed to be
// simple and easy to use.  It is not intended to be a general-purpose
// configuration library with robust type checking and validation.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// In this package, a configuration is represented as a map[string]any and
// serialized to a file in JSON format.
type Config map[string]any

// Write writes a configuration to a file as a JSON string.  Write returns and
// error if the configuration cannot be serialized or if the file cannot be
// written.
func (c Config) Write(filepath string) (err error) {
	text, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Unable to serialize config: %v", err)
	}
	err = os.WriteFile(filepath, text, 0644)
	if err != nil {
		return fmt.Errorf("Unable to write config file: %v", err)
	}
	return
}

// Read reads a configuration from a file. If the read is successful and the
// file text can be unmarshalled, Read overwrites the target configuration.
// Otherwise, an error is returned and the target configuration is not affected.
// Note that the behavior of json.Unmarshal is such that Unmarshal stores one of
// these in interface values:
//
// bool for JSON booleans,
// float64 for JSON numbers,
// string for JSON strings,
// []any for JSON arrays,
// map[string]any, for JSON objects,
// nil for JSON null.
//
// Thus, you can't get an integer back from a JSON number and any code that
// needs an integer will have to do the conversion itself.
func (c *Config) Read(filepath string) (err error) {
	text, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("Unable to read config file: %v", err)
	}
	tmp := make(map[string]any)
	err = json.Unmarshal(text, &tmp)
	if err != nil {
		return fmt.Errorf("Unable to parse config file: %v", err)
	}
	*c = tmp
	return
}

// Update updates a target configuration with the values from another
// configuration (source). If a key exists only in the source configuration,
// It is added to targer. If a key exists in both configurations, the value
// from the source configuration is used. Keys that are only present in the
// target configuration are not affected.
func (c Config) Update(source Config) {
	for k, v := range source {
		c[k] = v
	}
}

// HasKey returns true if the key exists in the configuration.
func (c Config) HasKey(key string) bool {
	_, ok := c[key]
	return ok
}

// HasKeyNested returns true if the nested key exists in the configuration.
func (c Config) HasKeyNested(keys ...string) bool {
	var ok bool
	var nestedMap = c
	var value any
	for i, key := range keys {
		value, ok = nestedMap[key]
		if !ok {
			return false
		}
		if i < len(keys)-1 {
			nestedMap = value.(map[string]interface{})
		}
	}

	return true
}

// UniqueKeyMatchOf returns the unique key, if any, that matches the given
// shortcut, k.  If there is no match or more than one match, the function
// returns an empty string. The comparison is case-insensitive and ignores the
// characters in the ignore slice. Thus, frob would match FROB, fRoB, and
// frobish.  If the ignore slice contains '_', frob would also match _f_ro_b_.
func (c Config) UniqueKeyMatchOf(k string, ignore []rune) string {
	ignoreMap := make(map[rune]bool)
	for _, r := range ignore {
		ignoreMap[r] = true
	}

	filter := func(r rune) rune {
		if ignoreMap[r] {
			return -1
		}
		return unicode.ToLower(r)
	}

	k = strings.Map(filter, k)
	var match string
	for key := range c {
		filteredKey := strings.Map(filter, key)
		if strings.HasPrefix(filteredKey, k) {
			if match != "" {
				return "" // More than one match, return empty string
			}
			match = key
		}
	}

	return match
}
