package config

import (
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// newViperConfigFromReader returns a viper config for content provided by passed reader.
func newViperConfigFromReader(reader io.Reader) (Config, error) {

	viperConfig := viper.New()
	viperConfig.SetConfigType("yaml")
	if err := viperConfig.ReadConfig(reader); err != nil {
		return nil, err
	}
	return &ViperConfig{config: viperConfig}, nil
}

// AsIntPtr will return passed int value as pointer.
func AsIntPtr(v int) *int {
	return &v
}

// AsStringPtr return given value as pointer.
func AsStringPtr(v string) *string {
	return &v
}

// AsBoolPtr returns given value as pointer.
func AsBoolPtr(v bool) *bool {
	return &v
}

// AsDurationPtr returns given value as pointer.
func AsDurationPtr(v time.Duration) *time.Duration {
	return &v
}

// toDuration will try to convert passed config value to a duration.
// Examples: 1s, 4m, 2h. Default unit is second, so config value 3 will be returned as 3 seconds.
func toDuration(value string) *time.Duration {

	if !isValidDuration(value) {
		return nil
	}

	duration := durationForUnit(value)
	if numValue := extractNumbers(value); numValue != nil {
		if intValue, err := strconv.ParseInt(*numValue, 10, 64); err == nil {
			duration = time.Duration(intValue) * duration
		}
	}
	return &duration
}

func durationForUnit(durationAsString string) time.Duration {

	switch true {
	case strings.HasSuffix(durationAsString, "s"):
		return 1 * time.Second
	case strings.HasSuffix(durationAsString, "m"):
		return 1 * time.Minute
	case strings.HasSuffix(durationAsString, "h"):
		return 1 * time.Hour
	default:
		return 1 * time.Second
	}
}

// extractNumbers try to get numbers from given config value.
func extractNumbers(strValue string) *string {

	intRegexp := regexp.MustCompile("[0-9]+")
	if match := intRegexp.FindString(strValue); match != "" {
		return &match
	} else {
		return nil
	}
}

// isValidDuration if passed config values is composed by a number
// followed by a single char of s, m or h. A single int value is valid as well.
func isValidDuration(value string) bool {

	durationWithUnit := regexp.MustCompile("^[0-9]+[smh,0-9]{0,1}$")
	return durationWithUnit.MatchString(value)
}
