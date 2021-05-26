package config

import (
	"strings"
)

// StaticConfigSource uses a static config passed during creating this config source.
type StaticConfigSource struct {

	// Stativ config in YAML format.
	yamlConfig string
}

// NewStaticConfigSource returns source with given static config values.
func NewStaticConfigSource(yamlConfig string) ConfigSource {
	return &StaticConfigSource{yamlConfig: yamlConfig}
}

// Load static config. This will create a new ViperConfig with static config content.
func (source *StaticConfigSource) Load() (Config, error) {

	reader := strings.NewReader(source.yamlConfig)
	return newViperConfigFromReader(reader)
}
