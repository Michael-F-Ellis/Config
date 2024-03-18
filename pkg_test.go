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
func TestConfig_Update(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		update Config
		want   Config
	}{
		{
			name:   "Update with scalars",
			config: Config{"x": 1, "y": 2, "key1": map[string]interface{}{"nestedKey": "value1"}, "key2": []interface{}{1, 2, 3}},
			update: Config{"y": 3, "z": 4},
			want:   Config{"x": 1, "y": 3, "z": 4, "key1": map[string]interface{}{"nestedKey": "value1"}, "key2": []interface{}{1, 2, 3}},
		},
		{
			name:   "Add nested map",
			config: Config{"a": 1},
			update: Config{"key1": map[string]any{"k1": "newv1", "k3": "v3"}},
			want:   Config{"a": 1, "key1": map[string]any{"k1": "newv1", "k3": "v3"}},
		},
		{
			name:   "Update nested maps",
			config: Config{"key1": map[string]any{"k1": "v1", "k2": "v2"}},
			update: Config{"key1": map[string]any{"k1": "newv1", "k3": "v3"}},
			want:   Config{"key1": map[string]any{"k1": "newv1", "k2": "v2", "k3": "v3"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.Update(tt.update)
			if diff := deep.Equal(tt.config, tt.want); diff != nil {
				t.Errorf("Config.Update() = %v, want %v", diff, nil)
			}
		})
	}
}
func TestConfig_Get(t *testing.T) {
	cfg := Config{
		"key1": "value1",
		"key2": map[string]interface{}{
			"nestedKey": "value2",
		},
		"key3": []interface{}{1, 2, 3},
	}

	tests := []struct {
		name   string
		keys   []string
		want   any
		exists bool
	}{
		{
			name:   "Get scalar",
			keys:   []string{"key1"},
			want:   "value1",
			exists: true,
		},
		{
			name:   "Get nested key",
			keys:   []string{"key2", "nestedKey"},
			want:   "value2",
			exists: true,
		},
		{
			name:   "Get array",
			keys:   []string{"key3"},
			want:   []interface{}{1, 2, 3},
			exists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := cfg.Get(tt.keys...)
			if ok != tt.exists {
				t.Errorf("Config.Get(), exists =%v, want %v", got, tt.exists)
			}
			if ok {
				if diff := deep.Equal(got, tt.want); diff != nil {
					t.Errorf("Config.Get() = %v", diff)
				}
			}
		})
	}
}
