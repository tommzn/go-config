package s3_test

import (
	"os"
	"testing"

	s3config "github.com/tommzn/go-config/s3"
	"github.com/stretchr/testify/suite"
)

type S3ConfigTestSuite struct {
	suite.Suite
}

func TestS3ConfigTestSuite(t *testing.T) {
	suite.Run(t, new(S3ConfigTestSuite))
}

func (suite *S3ConfigTestSuite) TestS3ConfigSource() {

	awsRegion1 := "eu-central-1"
	configSource1, err := s3config.NewS3ConfigSource("no-bucket", "no-key", &awsRegion1)
	suite.Nil(err)
	config1, err1 := configSource1.Load()
	suite.NotNil(err1)
	suite.Nil(config1)

	configSource2, err := s3config.NewS3ConfigSourceFromEnv()
	if configSource2 == nil {
		suite.T().Skip("Skip S3 tests. Missing env vars GO_CONFIG_S3_BUCKET and GO_CONFIG_S3_KEY to access S3 for test.")
	}
	suite.Nil(err)
	suite.NotNil(configSource2)

	// Exercise the actual download/parse path when integration env vars are configured,
	// instead of only checking that the source was constructed.
	loadedConfig, loadErr := configSource2.Load()
	suite.NoError(loadErr)
	suite.NotNil(loadedConfig)
}

func (suite *S3ConfigTestSuite) TestS3ConfigSourceWithoutRegion() {

	prev, hasPrev := os.LookupEnv("AWS_REGION")
	os.Unsetenv("AWS_REGION")
	suite.T().Cleanup(func() {
		if hasPrev {
			os.Setenv("AWS_REGION", prev)
		} else {
			os.Unsetenv("AWS_REGION")
		}
	})

	configSource, err := s3config.NewS3ConfigSource("no-bucket", "no-key", nil)
	suite.Nil(err)
	suite.NotNil(configSource)

	cfg, err := configSource.Load()
	suite.NotNil(err)
	suite.Nil(cfg)
}

func (suite *S3ConfigTestSuite) TestS3ConfigSourceFromEnvMissingBucket() {

	suite.T().Setenv("AWS_REGION", "eu-central-1")
	prevBucket, hasBucket := os.LookupEnv("GO_CONFIG_S3_BUCKET")
	os.Unsetenv("GO_CONFIG_S3_BUCKET")
	suite.T().Cleanup(func() {
		if hasBucket {
			os.Setenv("GO_CONFIG_S3_BUCKET", prevBucket)
		} else {
			os.Unsetenv("GO_CONFIG_S3_BUCKET")
		}
	})
	prevKey, hasKey := os.LookupEnv("GO_CONFIG_S3_KEY")
	os.Unsetenv("GO_CONFIG_S3_KEY")
	suite.T().Cleanup(func() {
		if hasKey {
			os.Setenv("GO_CONFIG_S3_KEY", prevKey)
		} else {
			os.Unsetenv("GO_CONFIG_S3_KEY")
		}
	})

	source, err := s3config.NewS3ConfigSourceFromEnv()
	suite.NotNil(err)
	suite.Nil(source)
	suite.Contains(err.Error(), "GO_CONFIG_S3_BUCKET")
}

func (suite *S3ConfigTestSuite) TestS3ConfigSourceFromEnvMissingKey() {

	suite.T().Setenv("AWS_REGION", "eu-central-1")
	suite.T().Setenv("GO_CONFIG_S3_BUCKET", "some-bucket")
	prevKey, hasKey := os.LookupEnv("GO_CONFIG_S3_KEY")
	os.Unsetenv("GO_CONFIG_S3_KEY")
	suite.T().Cleanup(func() {
		if hasKey {
			os.Setenv("GO_CONFIG_S3_KEY", prevKey)
		} else {
			os.Unsetenv("GO_CONFIG_S3_KEY")
		}
	})

	source, err := s3config.NewS3ConfigSourceFromEnv()
	suite.NotNil(err)
	suite.Nil(source)
	suite.Contains(err.Error(), "GO_CONFIG_S3_KEY")
}
