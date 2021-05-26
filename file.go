package config

import (
	"bytes"
	"io/ioutil"

	"github.com/spf13/viper"
)

// FileConfigSource reads a config file in YAML format using viper config.
type FileConfigSource struct {
	configFile *string
}

// FileConfigSource returns a new config source for given file.
// If you don't passed a specific config file this source will have a lool
// at different places for a default config file.
// See Load method for more details.
func NewFileConfigSource(configFile *string) ConfigSource {
	return &FileConfigSource{configFile: configFile}
}

// Load reads a config file and returns a ViperConfig.
// It uses the config file you've set during creating this source or
// it tries to find a file names config.yml or testconfig.yml in following locations.
// - loca directory, "./"
// - user home, "$HOMR/"
// - user home at go_config dir, "$HOMR/go_config/"
// - at "/etc/go_config/"
func (source *FileConfigSource) Load() (Config, error) {

	if source.configFile != nil {

		fileContent, err := ioutil.ReadFile(*source.configFile)
		if err != nil {
			return nil, err
		}
		return newViperConfigFromReader(bytes.NewReader(fileContent))
	}

	viperConfig := viper.New()
	viperConfig.AddConfigPath(".")
	viperConfig.AddConfigPath("$HOMR/")
	viperConfig.AddConfigPath("$HOMR/go_config/")
	viperConfig.AddConfigPath("/etc/go_config/")
	viperConfig.SetConfigName("config")
	viperConfig.SetConfigName("testconfig")
	viperConfig.SetConfigType("yaml")
	err := viperConfig.ReadInConfig()

	var config Config
	if err == nil {
		config = &ViperConfig{config: viperConfig}
	}
	return config, nil
}
