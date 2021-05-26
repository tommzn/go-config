package config

import (
	"time"
)

// ConfigSource can be used to load a config from different sources and
// with different formats.
type ConfigSource interface {

	// Load config depending on config loader implemanteation and
	// return it as Config.
	Load() (Config, error)
}

// Config is the ain interface provides by this package to get an access
// point for config from different sources and formats.
type Config interface {

	// Get try to load config value for passed key and will return given default
	// if it's not available.
	Get(key string, defaultValue *string) *string

	// GetAsInt try to load value for given config and will try to convert it
	// to imt. If there's no config value for passed key or conversion to int failes,
	// it wll return given default value.
	GetAsInt(key string, defaultValue *int) *int

	// GetAsIntSlice returns a string slice of config values for passed key
	// or return passed default value it there's no value for tis key.
	GetAsIntSlice(key string, defaultValue *[]int) *[]int

	// GetAsBool returns config value as bool or given default value
	// if there's no value for this key or conversion to bool fails.
	GetAsBool(key string, defaultValue *bool) *bool

	// GetAsDuration returns config value as duration or passed default value
	// if there's no value for passed key or maybe config value parsing to duration fails.
	// Unit for durations can defined with suffix "s" for seconds, "m" for minutes or "h" for hourse.
	// If there's no unit default will be seconds.
	GetAsDuration(key string, defaultValue *time.Duration) *time.Duration

	// GetSliceOfMap returns all config values as a slice of maps.
	GetAsSliceOfMaps(key string) []map[string]string
}
