package config

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

// errReader is an io.Reader that always fails, for exercising read-error paths.
type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

type ParserTestSuite struct {
	suite.Suite
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

func (suite *ParserTestSuite) TestParseConfigInvalidJSON() {

	defer SetConfigType("yaml")
	SetConfigType("json")

	config, err := parseConfig([]byte("{not valid json"))
	suite.NotNil(err)
	suite.Nil(config)
}

func (suite *ParserTestSuite) TestParseConfigInvalidYAML() {

	config, err := parseConfig([]byte("key: [unterminated"))
	suite.NotNil(err)
	suite.Nil(config)
}

func (suite *ParserTestSuite) TestLookupOnNonMapIntermediate() {

	conf := &mapConfig{values: map[string]interface{}{"key2": "value2"}}

	value, ok := conf.lookup("key2.nested")
	suite.False(ok)
	suite.Nil(value)
}

func (suite *ParserTestSuite) TestGetAsBoolInvalidValue() {

	conf := &mapConfig{values: map[string]interface{}{"flag": "not-a-bool"}}

	defaultValue := AsBoolPtr(true)
	value := conf.GetAsBool("flag", defaultValue)
	suite.Equal(defaultValue, value)
}

func (suite *ParserTestSuite) TestGetAsIntSliceOnNonSliceValue() {

	conf := &mapConfig{values: map[string]interface{}{"key2": "value2"}}

	defaultValue := &[]int{1, 2}
	value := conf.GetAsIntSlice("key2", defaultValue)
	suite.Equal(defaultValue, value)
}

func (suite *ParserTestSuite) TestToStringTypes() {

	cases := []struct {
		name  string
		input interface{}
		want  string
	}{
		{"int", 42, "42"},
		{"int64", int64(42), "42"},
		{"float64", 3.14, "3.14"},
		{"bool", true, "true"},
		{"unsupported-slice", []interface{}{1, 2}, ""},
	}
	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.Equal(tc.want, toString(tc.input))
		})
	}
}

func (suite *ParserTestSuite) TestToIntTypes() {

	cases := []struct {
		name  string
		input interface{}
		want  int
	}{
		{"int", 42, 42},
		{"int64", int64(42), 42},
		{"float64", float64(42), 42},
		{"numeric-string", "42", 42},
		{"non-numeric-string", "not-a-number", 0},
		{"unsupported-bool", true, 0},
	}
	for _, tc := range cases {
		suite.Run(tc.name, func() {
			suite.Equal(tc.want, toInt(tc.input))
		})
	}
}

func (suite *ParserTestSuite) TestNewConfigFromReader() {

	config, err := NewConfigFromReader(strings.NewReader("key: val"))
	suite.Nil(err)
	suite.NotNil(config)

	value := config.Get("key", nil)
	suite.NotNil(value)
	suite.Equal("val", *value)
}

func (suite *ParserTestSuite) TestParseConfigFromReaderReadError() {

	config, err := parseConfigFromReader(errReader{})
	suite.NotNil(err)
	suite.Nil(config)
}

func (suite *ParserTestSuite) TestGetAsSliceOfMapsSkipsNonMapAndEmptyItems() {

	conf := &mapConfig{values: map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"a": "val"},
			"not-a-map",
			map[string]interface{}{"b": 123}, // no string values -> toStringMap is empty, skipped
		},
	}}

	result := conf.GetAsSliceOfMaps("items")
	suite.Len(result, 1)
	suite.Equal("val", result[0]["a"])
}
