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

	suite.Equal("42", toString(42))
	suite.Equal("42", toString(int64(42)))
	suite.Equal("3.14", toString(3.14))
	suite.Equal("true", toString(true))
	suite.Equal("", toString([]interface{}{1, 2}))
}

func (suite *ParserTestSuite) TestToIntTypes() {

	suite.Equal(42, toInt(42))
	suite.Equal(42, toInt(int64(42)))
	suite.Equal(42, toInt(float64(42)))
	suite.Equal(42, toInt("42"))
	suite.Equal(0, toInt("not-a-number"))
	suite.Equal(0, toInt(true))
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
