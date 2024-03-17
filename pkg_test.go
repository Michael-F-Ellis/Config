package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var tempDir string

func TestMain(m *testing.M) {
	// Setup
	var err error
	tempDir = os.TempDir()
	if err != nil {
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Teardown
	os.RemoveAll(tempDir)

	os.Exit(code)
}
func TestWriteAndRead(t *testing.T) {
	// Setup
	cfg := Config{ /* initialize with desired values */ }
	filePath := filepath.Join(tempDir, "test_config.json")
	defer os.Remove(filePath)

	// Test Write
	err := Write(filePath, cfg)
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	// Test Read
	readCfg := Config{}
	err = Read(filePath, readCfg)
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}

	// Check if the read configuration matches the written one
	if !reflect.DeepEqual(cfg, readCfg) {
		t.Errorf("Read() = %v, want %v", readCfg, cfg)
	}
}
