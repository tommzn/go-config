package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileConfigSource reads a config file in YAML or JSON format (see SetConfigType).
type FileConfigSource struct {
	configFile *string
}

// NewFileConfigSource returns a new config source for given file.
// If you don't passed a specific config file this source will have a lool
// at different places for a default config file.
// See Load method for more details.
func NewFileConfigSource(configFile *string) ConfigSource {
	return &FileConfigSource{configFile: configFile}
}

// Load reads a config file and returns a Config.
// It uses the config file you've set during creating this source or, if none was
// given, looks for a file named config.yaml/config.yml (or config.json, if the
// config type is set to "json") in following locations, in order:
// - local directory, "./"
// - user home, "$HOME/"
// - user home at go_config dir, "$HOME/go_config/"
// - at "/etc/go_config/"
func (source *FileConfigSource) Load() (Config, error) {

	if source.configFile != nil {
		fileContent, err := os.ReadFile(*source.configFile)
		if err != nil {
			return nil, fmt.Errorf("reading config file %s: %w", *source.configFile, err)
		}
		return parseConfig(fileContent)
	}

	searchPaths := []string{".", os.ExpandEnv("$HOME"), os.ExpandEnv("$HOME/go_config"), "/etc/go_config"}
	for _, dir := range searchPaths {
		for _, name := range defaultConfigFileNames() {
			if content, err := os.ReadFile(filepath.Join(dir, name)); err == nil {
				return parseConfig(content)
			}
		}
	}
	return nil, fmt.Errorf("no config file found in %s", strings.Join(searchPaths, ", "))
}

// defaultConfigFileNames returns the file names Load will look for when no
// explicit config file was given, depending on the currently configured config type.
func defaultConfigFileNames() []string {
	if strings.EqualFold(configType, "json") {
		return []string{"config.json"}
	}
	return []string{"config.yaml", "config.yml"}
}
