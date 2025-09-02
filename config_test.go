package config

import (
	"io/ioutil"
	"time"

	"github.com/stretchr/testify/suite"
	//"log"

	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) TestDefaultConfigSource() {

	configSource := NewConfigSource()
	suite.testConfigSource(configSource)
}

func (suite *ConfigTestSuite) TestStaticConfigSource() {

	configSource := NewStaticConfigSource(suite.staticConfigForTest())
	suite.testConfigSource(configSource)
}

func (suite *ConfigTestSuite) TestS3ConfigSource() {

	awsRegion1 := "eu-central-1"
	configSource1, err := NewS3ConfigSource("no-bucket", "no-key", &awsRegion1)
	suite.Nil(err)
	config1, err1 := configSource1.Load()
	suite.NotNil(err1)
	suite.Nil(config1)

	configSource2, err := NewS3ConfigSourceFromEnv()
	if configSource2 == nil {
		suite.T().Skip("Skip S3 tests. Missing env vars GO_CONFIG_S3_BUCKET and GO_CONFIG_S3_KEY to access S3 for test.")
	}
	suite.Nil(err)
	suite.testConfigSource(configSource2)
}

func (suite *ConfigTestSuite) TestFileConfigSource() {

	configSource1 := NewFileConfigSource(nil)
	suite.testConfigSource(configSource1)

	configFile2 := "./testconfig.yml"
	configSource2 := NewFileConfigSource(&configFile2)
	suite.testConfigSource(configSource2)

	configFile3 := "./notexistingfile.yml"
	configSource3 := NewFileConfigSource(&configFile3)
	config, err := configSource3.Load()
	suite.NotNil(err)
	suite.Nil(config)
}

func (suite *ConfigTestSuite) testConfigSource(configSource ConfigSource) {

	config, err := configSource.Load()
	suite.Nil(err)

	suite.testGetConfigValuesAsString(config)
	suite.testGetStructuredConfigValue(config)
	suite.testGetConfigValuesAsInt(config)
	suite.testGetConfigValuesAsBool(config)
	suite.testGetConfigValuesAsIntSlice(config)
	suite.testGetConfigValuesAsSliceOfMaps(config)
	suite.testGetConfigValuesAsDuration(config)
}

func (suite *ConfigTestSuite) testGetConfigValuesAsString(config Config) {

	configKey := "key2"
	expectedValue := "value2"
	defaultValue := AsStringPtr("DefaultValue")
	notExistingKey := "xxx"

	value1 := config.Get(configKey, nil)
	suite.NotNil(value1)
	suite.Equal(expectedValue, *value1)

	value2 := config.Get(notExistingKey, defaultValue)
	suite.NotNil(value2)
	suite.Equal(*defaultValue, *value2)

	value3 := config.Get(notExistingKey, nil)
	suite.Nil(value3)

}

func (suite *ConfigTestSuite) testGetStructuredConfigValue(config Config) {

	configKey := "namespace1.key1"
	expectedValue := "value1"
	defaultValue := AsStringPtr("DefaultValue")
	notExistingKey := "xxx"

	value1 := config.Get(configKey, nil)
	suite.NotNil(value1)
	suite.Equal(expectedValue, *value1)

	value2 := config.Get(notExistingKey, defaultValue)
	suite.NotNil(value2)
	suite.Equal(*defaultValue, *value2)

	value3 := config.Get(notExistingKey, nil)
	suite.Nil(value3)

}

func (suite *ConfigTestSuite) testGetConfigValuesAsInt(config Config) {

	configKey := "key3"
	expectedValue := 12345
	defaultValue := AsIntPtr(6789)
	notExistingKey := "xxx"

	value1 := config.GetAsInt(configKey, nil)
	suite.NotNil(value1)
	suite.Equal(expectedValue, *value1)

	value2 := config.GetAsInt(notExistingKey, defaultValue)
	suite.NotNil(value2)
	suite.Equal(*defaultValue, *value2)

	value3 := config.GetAsInt(notExistingKey, nil)
	suite.Nil(value3)

}

func (suite *ConfigTestSuite) testGetConfigValuesAsBool(config Config) {

	configKey := "boolval"
	expectedValue := true
	defaultValue := AsBoolPtr(false)
	notExistingKey := "xxx"

	value1 := config.GetAsBool(configKey, nil)
	suite.NotNil(value1)
	suite.Equal(expectedValue, *value1)

	value2 := config.GetAsBool(notExistingKey, defaultValue)
	suite.NotNil(value2)
	suite.Equal(*defaultValue, *value2)

	value3 := config.GetAsBool(notExistingKey, nil)
	suite.Nil(value3)

}

func (suite *ConfigTestSuite) testGetConfigValuesAsIntSlice(config Config) {

	configKey := "intslice"
	expectedValue := []int{342543545, 3465567, 547657}
	defaultValue := &[]int{1, 2}
	notExistingKey := "xxx"

	value1 := config.GetAsIntSlice(configKey, nil)
	suite.NotNil(value1)
	suite.Equal(expectedValue, *value1)

	value2 := config.GetAsIntSlice(notExistingKey, defaultValue)
	suite.NotNil(value2)
	suite.Equal(*defaultValue, *value2)

	value3 := config.GetAsIntSlice(notExistingKey, nil)
	suite.Nil(value3)

	value4 := config.GetAsInt("intslice2", nil)
	suite.NotNil(value4)
	suite.Equal(0, *value4)

}

func (suite *ConfigTestSuite) testGetConfigValuesAsSliceOfMaps(config Config) {

	configKey := "sliceofmaps"
	expectedSize := 2
	notExistingKey := "xxx"

	value1 := config.GetAsSliceOfMaps(configKey)
	suite.NotNil(value1)
	suite.Len(value1, expectedSize)

	value2 := config.GetAsSliceOfMaps(notExistingKey)
	suite.Len(value2, 0)
}

func (suite *ConfigTestSuite) testGetConfigValuesAsDuration(config Config) {

	duration1 := config.GetAsDuration("durations.seconds", nil)
	suite.NotNil(duration1)
	suite.Equal(43*time.Second, *duration1)

	duration2 := config.GetAsDuration("durations.minutes", nil)
	suite.NotNil(duration2)
	suite.Equal(21*time.Minute, *duration2)

	duration3 := config.GetAsDuration("durations.hours", nil)
	suite.NotNil(duration3)
	suite.Equal(5*time.Hour, *duration3)

	duration4 := config.GetAsDuration("durations.defaultvalue", nil)
	suite.NotNil(duration4)
	suite.Equal(22*time.Second, *duration4)

	duration5 := config.GetAsDuration("durations.unsupported", nil)
	suite.Nil(duration5)

	duration6 := config.GetAsDuration("durations.notexisting", nil)
	suite.Nil(duration6)
}

// staticConfigForTest returns a static config in YAML format.
func (suite *ConfigTestSuite) staticConfigForTest() string {
	fileContent, err := ioutil.ReadFile("testconfig.yml")
	suite.Nil(err)
	return string(fileContent)
}
