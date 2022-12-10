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

func (suite *UtilsTestSuite) TestNewViperConfigFromReader() {

	configStr1 := "key: val"
	config1, err1 := newViperConfigFromReader(strings.NewReader(configStr1))
	suite.Nil(err1)
	suite.NotNil(config1)

	configStr2 := "key1=val1"
	config2, err2 := newViperConfigFromReader(strings.NewReader(configStr2))
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

	duration1 := toDuration("7s")
	suite.NotNil(duration1)
	suite.Equal(7*time.Second, *duration1)

	duration2 := toDuration("5m")
	suite.NotNil(duration2)
	suite.Equal(5*time.Minute, *duration2)

	duration3 := toDuration("2h")
	suite.NotNil(duration3)
	suite.Equal(2*time.Hour, *duration3)

	duration4 := toDuration("11")
	suite.NotNil(duration4)
	suite.Equal(11*time.Second, *duration4)

	duration5 := AsDuration("7s")
	suite.NotNil(duration5)
	suite.Equal(7*time.Second, *duration5)

	suite.Nil(toDuration("3d"))
	suite.Nil(toDuration("xxx"))
	suite.Nil(toDuration("ABCs"))
}

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
}

func (suite *UtilsTestSuite) TestSetViperConfigType() {

	suite.Equal("yaml", viperConfigType)

	jsonConfigType := "json"
	SetViperConfigType(jsonConfigType)
	suite.Equal(jsonConfigType, viperConfigType)
}
