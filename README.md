# Config
A Go package for creating and using JSON configuration files.

Package `config` provides methods to read, write and update configurations.

The intended use case is in programs that make calls to API's that accept POST requests with JSON payloads.  

The package is designed to be simple and easy to use.  It is not intended to be a general-purpose configuration library with robust type checking and validation and support for multiple formats.

On the other hand, JSON is widely used and the approach taken in the design of `config` can cope with nested JSON objects to an arbitrary depth.

## Install
```bash
$> go get github.com/Michael-F-Ellis/config
```
## Usage:
```go
import github.com/Michael-F-Ellis/config
// Creating a Config
cfg := Config{
		/* initialize with desired values */
		"foo": "bar",
		"baz": 42.,
		"boo": true,
		"qux": []any{"quux", "corge"},
		"g_rault": map[string]any{
			"garply": "waldo",
			"fred":   42.42,
		}
// Serializing
err := cfg.Write("some/path/myconfig.json")

// Reading
othercfg, err := cfg.Read("some/other/path/otherconfig.json")

// Updating (in memory)
cfg.Update(othercfg) // replaces or adds key:value pairs

// Check for key existence
has := cfg.HasKey("baz") // true

// Check nested key existence
has = cfg.HasKeyNested("grault", "fred") // true

// Validate and use a shortcut
s := cfg.UniqueMatchOf("GR", '_') // g_rault matches
if s != "" {
	value := cfg[s]["fred"]
}
```
