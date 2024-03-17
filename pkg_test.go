package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-test/deep"
)

var tempDir string

func TestMain(m *testing.M) {
	// Setup
	var err error
	tempDir = os.TempDir()
	// check that the temp directory exists
	if _, err = os.Stat(tempDir); os.IsNotExist(err) {
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	os.Exit(code)
}
func TestWriteAndRead(t *testing.T) {
	// Setup
	cfg := Config{
		/* initialize with desired values */
		"foo": "bar",
		"baz": 42.,
		"boo": true,
		"qux": []any{"quux", "corge"},
		"grault": map[string]any{
			"garply": "waldo",
			"fred":   42.42,
		},
	}
	filePath := filepath.Join(tempDir, "test_config.json")
	defer os.Remove(filePath)

	// Test Write
	err := cfg.Write(filePath)
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	// Test Read
	readCfg := Config{}
	err = readCfg.Read(filePath)
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}

	// Check if the read configuration matches the written one
	if diff := deep.Equal(cfg, readCfg); diff != nil {
		t.Errorf("%v", diff)
	}
}
func TestHasKeyNested(t *testing.T) {
	cfg := Config{
		"key0": 42,
		"key1": map[string]interface{}{
			"key2": map[string]interface{}{
				"key3": "value",
			},
		},
	}

	tests := []struct {
		name string
		keys []string
		want bool
	}{
		{
			name: "Existing unnested key",
			keys: []string{"key0"},
			want: true,
		},
		{
			name: "Existing nested keys",
			keys: []string{"key1", "key2", "key3"},
			want: true,
		},

		{
			name: "Non-existing keys",
			keys: []string{"key1", "key2", "nonexistent"},
			want: false,
		},
		{
			name: "More non-existing keys",
			keys: []string{"keyx", "key1", "key2"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cfg.HasKeyNested(tt.keys...); got != tt.want {
				t.Errorf("HasKeyNested() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestUniqueKeyMatchOf(t *testing.T) {
	cfg := Config{
		"key1":    "value1",
		"ukey2":   "value2",
		"ukey3":   "value3",
		"_x_key_": "xvalue",
	}

	tests := []struct {
		name   string
		k      string
		ignore []rune
		want   string
	}{
		{
			name:   "Unique prefix",
			k:      "k",
			ignore: []rune{},
			want:   "key1",
		},
		{
			name:   "Non-unique prefix",
			k:      "uk",
			ignore: []rune{' '},
			want:   "",
		},
		{
			name:   "Ignore characters",
			k:      "xkey",
			ignore: []rune{'_'},
			want:   "_x_key_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cfg.UniqueKeyMatchOf(tt.k, tt.ignore); got != tt.want {
				t.Errorf("UniqueKeyMatchOf %s = %v, want %v", tt.k, got, tt.want)
			}
		})
	}
}
