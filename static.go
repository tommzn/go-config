package config

import (
	"strings"
)

// StaticConfigSource uses a static config passed during creating this config source.
type StaticConfigSource struct {

	// Static config content, in YAML or JSON format (see SetConfigType).
	content string
}

// NewStaticConfigSource returns source with given static config values.
func NewStaticConfigSource(content string) ConfigSource {
	return &StaticConfigSource{content: content}
}

// Load static config. This will parse the static content passed at creation.
func (source *StaticConfigSource) Load() (Config, error) {

	reader := strings.NewReader(source.content)
	return parseConfigFromReader(reader)
}
