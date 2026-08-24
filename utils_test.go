package config

import (
	"strings"
	"time"

	"github.com/stretchr/testify/suite"
	//"log"

	"testing"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (suite *UtilsTestSuite) TestParseConfigFromReader() {

	defer SetConfigType("yaml")
	SetConfigType("yaml")

	configStr1 := "key: val"
	config1, err1 := parseConfigFromReader(strings.NewReader(configStr1))
	suite.Nil(err1)
	suite.NotNil(config1)

	// Not a mapping at the top level - fails to unmarshal into map[string]interface{}.
	configStr2 := "key1=val1"
	config2, err2 := parseConfigFromReader(strings.NewReader(configStr2))
	suite.NotNil(err2)
	suite.Nil(config2)
}

func (suite *UtilsTestSuite) TestPointerConverter() {

	v1 := "Test"
	p1 := AsStringPtr(v1)
	suite.Equal(v1, *p1)

	v2 := 1234
	p2 := AsIntPtr(v2)
	suite.Equal(v2, *p2)

	v3 := true
	p3 := AsBoolPtr(v3)
	suite.Equal(v3, *p3)

	v4 := 2 * time.Second
	p4 := AsDurationPtr(v4)
	suite.Equal(v4, *p4)
}

func (suite *UtilsTestSuite) TestConvertToDuration() {

	cases := []struct {
		name  string
		input string
		want  *time.Duration
	}{
		{"seconds", "7s", ptrDuration(7 * time.Second)},
		{"minutes", "5m", ptrDuration(5 * time.Minute)},
		{"hours", "2h", ptrDuration(2 * time.Hour)},
		{"bare-number-defaults-to-seconds", "11", ptrDuration(11 * time.Second)},
		{"unknown-unit-days", "3d", nil},
		{"garbage", "xxx", nil},
		{"non-numeric-prefix", "ABCs", nil},
		{"trailing-comma", "12,", nil},
		{"trailing-dash", "12-", nil},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			got := toDuration(tc.input)
			if tc.want == nil {
				suite.Nil(got)
				return
			}
			suite.NotNil(got)
			suite.Equal(*tc.want, *got)
		})
	}

	suite.Run("as-duration-public-wrapper", func() {
		got := AsDuration("7s")
		suite.NotNil(got)
		suite.Equal(7*time.Second, *got)
	})
}

func ptrDuration(d time.Duration) *time.Duration { return &d }

func (suite *UtilsTestSuite) TestExtracNumbers() {

	numbersInString := "24s"
	numbers := extractNumbers(numbersInString)
	suite.NotNil(numbers)
	suite.Equal("24", *numbers)

	suite.Nil(extractNumbers("xxx"))
}

func (suite *UtilsTestSuite) TestIsValidDuration() {

	suite.True(isValidDuration("2s"))
	suite.True(isValidDuration("5s"))
	suite.True(isValidDuration("3h"))
	suite.True(isValidDuration("11"))
	suite.False(isValidDuration("1d"))
	suite.False(isValidDuration("8y"))
	suite.False(isValidDuration("ABC"))
	suite.False(isValidDuration("12,"))  // trailing comma must be rejected
	suite.False(isValidDuration("12-"))  // other non-unit trailing chars must be rejected
}

func (suite *UtilsTestSuite) TestSetConfigType() {

	defer SetConfigType("yaml")
	suite.Equal("yaml", configType)

	jsonConfigType := "json"
	SetConfigType(jsonConfigType)
	suite.Equal(jsonConfigType, configType)
}
