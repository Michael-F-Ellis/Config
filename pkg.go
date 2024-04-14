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
	"reflect"
	"strings"
	"unicode"
)

// In this package, a configuration is represented as a map[string]any and
// serialized to a file in JSON format.
type Config map[string]any

// ConfigFromString returns a config from a string
// The string should be in the form of a JSON object.
// Opening and closing braces are optional.
func ConfigFromString(s string) (map[string]any, error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "{") {
		s = "{" + s
	}
	if !strings.HasSuffix(s, "}") {
		s = s + "}"
	}
	var c = map[string]any{}
	err := json.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Write writes a configuration to a file as a JSON string.  Write returns and
// error if the configuration cannot be serialized or if the file cannot be
// written.
func Write(cfg map[string]any, filepath string) (err error) {
	text, err := json.Marshal(cfg)
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
func Read(filepath string) (c map[string]any, err error) {
	text, err := os.ReadFile(filepath)
	if err != nil {
		err = fmt.Errorf("Unable to read config file: %v", err)
		return
	}
	tmp := make(map[string]any)
	err = json.Unmarshal(text, &tmp)
	if err != nil {
		err = fmt.Errorf("Unable to parse config file: %v", err)
		return
	}
	c = tmp
	return
}

// Update updates a target configuration with the values from another
// configuration (source). If a key exists only in the source configuration,
// It is added to target. If a key exists in both configurations, the value
// from the source configuration is used. Keys that are only present in the
// target configuration are not affected. Map values are updated recursively.
func Update(source, target map[string]any) {
	for k, v := range source {
		switch v := v.(type) {
		case map[string]any:
			if _, ok := target[k]; ok {
				// If the key exists in both configurations, update the target
				// recursively.
				tmp := target[k].(map[string]any)
				Update(v, tmp)
			} else {
				// If the key only exists in the source configuration, add it to
				// the target
				target[k] = v
			}
		// If the value is not a map, i.e. a number, string, bool, null or array
		// just add it to the target configuration, replacing the value if the
		// key already exists.
		default:
			target[k] = v
		}
	}
}

// HasKey returns true if the key exists in the configuration.
func HasKey(c map[string]any, key string) bool {
	_, ok := c[key]
	return ok
}

// HasKeyNested returns true if the nested key exists in the configuration.
func HasKeyNested(c map[string]any, keys ...string) bool {
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

// Get returns the value of a (nested) key in the configuration and a boolean
// indicating whether the key exists.
func Get(c map[string]any, keys ...string) (value any, ok bool) {
	nestedMap := c
	for i, key := range keys {
		value, ok = nestedMap[key]
		if !ok {
			return
		}
		if i < len(keys)-1 {
			nestedMap = value.(map[string]interface{})
		}
	}
	return
}

// Set sets the value of a (nested) key in the configuration. If the key does
// not exist, it is added to the configuration. If the key exists, the value is
// updated. If the key is nested, the nested map is created if it does not
// exist.
func Set(c map[string]any, value any, keys ...string) {
	nestedMap := c
	for i, key := range keys {
		if i < len(keys)-1 {
			if _, ok := nestedMap[key]; !ok {
				nestedMap[key] = map[string]any{}
			}
			nestedMap = nestedMap[key].(map[string]any)
		} else {
			nestedMap[key] = value
		}
	}
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

type Translation map[string]string

// Apply() interprets the keys and values of a Translation object
// as paths to nested values in a Config object. It splits the keys
// and values by the separator and then sets the values in cTo
// to the corresponding values in cFrom.
func (t Translation) Apply(cFrom, cTo map[string]any, sep string) error {
	for k, v := range t {
		toKeys := strings.Split(v, sep)
		for i := range toKeys {
			toKeys[i] = strings.TrimSpace(toKeys[i])
		}
		fromKeys := strings.Split(k, sep)
		for i := range fromKeys {
			fromKeys[i] = strings.TrimSpace(fromKeys[i])
		}
		fromV, ok := Get(cFrom, fromKeys...)
		if !ok {
			return fmt.Errorf("key %v not found in cFrom", fromKeys)
		}
		Set(cTo, fromV, toKeys...)

	}
	return nil
}

// CompareTypes compares the types of keys in two Config objects.  The function
// compares the types of keys in the Config object with the keys in the reference
// Config object.  The function returns two slices of strings.  The first slice
// contains the keys where the types do not match.  The second slice contains the
// keys that are in the Config object but not in the reference Config object.
func (c Config) CompareTypes(refCfg Config, prefix string, mismatches *[]string, notFound *[]string) {
	for key, value1 := range c {
		// build prefix for nested keys
		pfx := strings.Join([]string{prefix, key}, ":")
		pfx = strings.TrimLeft(pfx, ":") // remove leading colon caused by empty prefix
		// check if key exists in refCfg
		if value2, ok := refCfg[key]; ok {
			type1 := reflect.TypeOf(value1)
			type2 := reflect.TypeOf(value2)
			if type1 == type2 {
				if reflect.ValueOf(value1).Kind() == reflect.Map {
					// recurse into nested map
					ck := value1.(Config)
					ck.CompareTypes(value2.(Config), pfx, mismatches, notFound)
				}
			} else {
				// found a type mismatch
				*mismatches = append(*mismatches, fmt.Sprintf("%s:%s!=%s", pfx, type1, type2))
			}
		} else {
			// key not found in map2
			*notFound = append(*notFound, pfx)
		}
	}
}
