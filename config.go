// Package config provides access to config from different sources in YAML format.
// Uses viper config from github.com/spf13/viper to load and access config values.
package config

// NewConfigSource returns the default config loader, the ViperConfigSource.
func NewConfigSource() ConfigSource {
	return NewFileConfigSource(nil)
}
