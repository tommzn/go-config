package config

import (
	"strconv"
	"time"

	"github.com/spf13/viper"
)

var viperConfigType string

// init sets type for viper config to YAML.
func init() {
	viperConfigType = "yaml"
}

// SetViperConfigType sets the config type which should be used at creating new viper config.
func SetViperConfigType(configType string) {
	viperConfigType = configType
}

// ViperConfig is a wrapper to config handler provided by github.com/spf13/viper.
type ViperConfig struct {
	config *viper.Viper
}

// Get try to load config value for passed key and will return given default
// if it's not available.
func (conf *ViperConfig) Get(key string, defaultValue *string) *string {
	if conf.config.IsSet(key) {
		value := conf.config.GetString(key)
		return &value
	}
	return defaultValue
}

// GetAsInt try to load value for given config and will try to convert it
// to imt. If there's no config value for passed key or conversion to int failes,
// it wll return given default value.
func (conf *ViperConfig) GetAsInt(key string, defaultValue *int) *int {
	if conf.config.IsSet(key) {
		value := conf.config.GetInt(key)
		return &value
	}
	return defaultValue
}

// GetAsIntSlice returns a string slice of config values for passed key
// or return passed default value it there's no value for tis key.
func (conf *ViperConfig) GetAsIntSlice(key string, defaultValue *[]int) *[]int {
	if conf.config.IsSet(key) {
		value := conf.config.GetIntSlice(key)
		return &value
	}
	return defaultValue
}

// GetAsBool returns config value as bool or given default value
// if there's no value for this key or conversion to bool fails.
func (conf *ViperConfig) GetAsBool(key string, defaultValue *bool) *bool {
	if conf.config.IsSet(key) {
		value := conf.config.GetString(key)
		b, err := strconv.ParseBool(value)
		if err == nil {
			return &b
		}
	}
	return defaultValue
}

// GetAsDuration returns config value as duration or passed default value
// if there's no value for passed key or maybe config value parsing to duration fails.
// Unit for durations can defined with suffix "s" for seconds, "m" for minutes or "h" for hourse.
// If there's no unit default will be seconds.
func (conf *ViperConfig) GetAsDuration(key string, defaultValue *time.Duration) *time.Duration {

	if conf.config.IsSet(key) {
		value := conf.config.GetString(key)
		return toDuration(value)
	}
	return defaultValue
}

// GetAsSliceOfMaps returns local config values as slice of maps.
func (conf *ViperConfig) GetAsSliceOfMaps(key string) []map[string]string {

	var retValues []map[string]string

	configValue := conf.config.Get(key)
	if configValue == nil {
		return retValues
	}
	if configSlice, ok := configValue.([]interface{}); ok {

		for _, configItem := range configSlice {

			if configMap, ok := configItem.(map[interface{}]interface{}); ok {

				stringMap := conf.toStringMap(configMap)
				if len(stringMap) > 0 {
					retValues = append(retValues, stringMap)
				}
			}
		}
	}
	return retValues
}

// toStringMap try to convert passed map with interface values to a map with string keys and values.
func (conf *ViperConfig) toStringMap(interfaceMap map[interface{}]interface{}) map[string]string {

	stringMap := make(map[string]string)
	for key, val := range interfaceMap {
		if strKey, okKey := key.(string); okKey {
			if strVal, okVal := val.(string); okVal {
				stringMap[strKey] = strVal
			}
		}
	}
	return stringMap
}
