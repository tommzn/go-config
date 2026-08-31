package config

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	mapstructure "github.com/go-viper/mapstructure/v2"
	"gopkg.in/yaml.v3"
)

// configType selects the format used to parse config content passed to
// NewFileConfigSource, NewStaticConfigSource and NewConfigFromReader.
// Defaults to "yaml". Change with SetConfigType.
var configType string

// init sets the default config type to YAML.
func init() {
	configType = "yaml"
}

// SetConfigType sets the format used to parse config content.
// Supported values: "yaml" (default), "json".
func SetConfigType(newConfigType string) {
	configType = newConfigType
}

// mapConfig is a Config implementation backed by a parsed, nested map of
// config values - the result of unmarshaling YAML or JSON content.
type mapConfig struct {
	values map[string]interface{}
}

// parseConfig parses passed content as YAML or JSON, depending on the
// currently configured config type, and returns a Config.
func parseConfig(content []byte) (Config, error) {

	var values map[string]interface{}
	var err error
	if strings.EqualFold(configType, "json") {
		err = json.Unmarshal(content, &values)
	} else {
		err = yaml.Unmarshal(content, &values)
	}
	if err != nil {
		return nil, err
	}
	return &mapConfig{values: lowercaseKeys(values)}, nil
}

// parseConfigFromReader reads all content from passed reader and parses it. See parseConfig.
func parseConfigFromReader(reader io.Reader) (Config, error) {

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	return parseConfig(content)
}

// lowercaseKeys recursively lowercases all string map keys of passed map, so
// config key lookups are case-insensitive.
func lowercaseKeys(values map[string]interface{}) map[string]interface{} {

	lowered := make(map[string]interface{}, len(values))
	for key, value := range values {
		if nested, ok := value.(map[string]interface{}); ok {
			value = lowercaseKeys(nested)
		}
		lowered[strings.ToLower(key)] = value
	}
	return lowered
}

// lookup walks passed dot-separated key through the nested config values and
// returns the value found and true, or nil and false if there's no value for it.
func (conf *mapConfig) lookup(key string) (interface{}, bool) {

	current := interface{}(conf.values)
	for _, part := range strings.Split(strings.ToLower(key), ".") {
		values, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current, ok = values[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

// Get try to load config value for passed key and will return given default
// if it's not available.
func (conf *mapConfig) Get(key string, defaultValue *string) *string {

	value, ok := conf.lookup(key)
	if !ok {
		return defaultValue
	}
	strValue := toString(value)
	return &strValue
}

// GetAsInt try to load value for given config and will try to convert it
// to imt. If there's no config value for passed key or conversion to int failes,
// it wll return given default value.
func (conf *mapConfig) GetAsInt(key string, defaultValue *int) *int {

	value, ok := conf.lookup(key)
	if !ok {
		return defaultValue
	}
	intValue := toInt(value)
	return &intValue
}

// GetAsIntSlice returns a string slice of config values for passed key
// or return passed default value it there's no value for tis key.
func (conf *mapConfig) GetAsIntSlice(key string, defaultValue *[]int) *[]int {

	value, ok := conf.lookup(key)
	if !ok {
		return defaultValue
	}
	items, ok := value.([]interface{})
	if !ok {
		return defaultValue
	}
	intSlice := make([]int, len(items))
	for i, item := range items {
		intSlice[i] = toInt(item)
	}
	return &intSlice
}

// GetAsBool returns config value as bool or given default value
// if there's no value for this key or conversion to bool fails.
func (conf *mapConfig) GetAsBool(key string, defaultValue *bool) *bool {

	value, ok := conf.lookup(key)
	if !ok {
		return defaultValue
	}
	if boolValue, err := strconv.ParseBool(toString(value)); err == nil {
		return &boolValue
	}
	return defaultValue
}

// GetAsDuration returns config value as duration or passed default value
// if there's no value for passed key or maybe config value parsing to duration fails.
// Unit for durations can defined with suffix "s" for seconds, "m" for minutes or "h" for hourse.
// If there's no unit default will be seconds.
func (conf *mapConfig) GetAsDuration(key string, defaultValue *time.Duration) *time.Duration {

	value, ok := conf.lookup(key)
	if !ok {
		return defaultValue
	}
	if duration := toDuration(toString(value)); duration != nil {
		return duration
	}
	return defaultValue
}

// GetAsSliceOfMaps returns local config values as slice of maps.
func (conf *mapConfig) GetAsSliceOfMaps(key string) []map[string]string {

	var retValues []map[string]string

	value, ok := conf.lookup(key)
	if !ok {
		return retValues
	}
	configSlice, ok := value.([]interface{})
	if !ok {
		return retValues
	}
	for _, configItem := range configSlice {
		if configMap, ok := configItem.(map[string]interface{}); ok {
			stringMap := toStringMap(configMap)
			if len(stringMap) > 0 {
				retValues = append(retValues, stringMap)
			}
		}
	}
	return retValues
}

// toStringMap try to convert passed map with interface values to a map with string keys and values.
func toStringMap(interfaceMap map[string]interface{}) map[string]string {

	stringMap := make(map[string]string)
	for key, val := range interfaceMap {
		if strVal, okVal := val.(string); okVal {
			stringMap[key] = strVal
		}
	}
	return stringMap
}

// Unmarshal decodes the configuration into the provided struct or map.
// The `rawVal` parameter should be a pointer to a struct or map where the
// configuration values will be unmarshaled. Returns an error if unmarshaling fails.
// Struct fields are matched using "mapstructure" tags, same as before.
func (conf *mapConfig) Unmarshal(rawVal any) error {
	return mapstructure.Decode(conf.values, rawVal)
}

// toString converts a decoded YAML/JSON scalar value to its string representation.
func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

// toInt converts a decoded YAML/JSON value to int, returning 0 if it can't be
// converted (e.g. the value is a slice or map) - matching the previous
// viper-backed implementation's behavior of returning the zero value rather
// than an error or the caller's default in that case.
func toInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if intValue, err := strconv.Atoi(v); err == nil {
			return intValue
		}
	}
	return 0
}
