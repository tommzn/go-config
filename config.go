// Package config provides access to config from different sources in YAML or JSON format.
// Config values can be loaded from files, static strings, or custom readers.
package config

// NewConfigSource returns the default config loader, a FileConfigSource that
// searches standard locations for a config file.
func NewConfigSource() ConfigSource {
	return NewFileConfigSource(nil)
}
