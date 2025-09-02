package config

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

// S3ConfigSource loads a YAML config from a file in an AWS S3 bucket.
type S3ConfigSource struct {

	// AWS config for S3 access.
	cfg aws.Config

	// Bucket the config file is located in.
	bucket string

	// Path and file name for a config file.
	key string
}

// NewS3ConfigSource returns a new S3 config source which uses the config file from the given S3 bucket.
// If region is empty it will try to get current AWS region from environment variable AWS_REGION.
func NewS3ConfigSource(bucket, key string, region *string) (ConfigSource, error) {
	var cfg aws.Config
	var err error

	if region != nil {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(*region))
	} else if envRegion, ok := os.LookupEnv("AWS_REGION"); ok {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(envRegion))
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		return nil, err
	}

	return &S3ConfigSource{
		cfg:    cfg,
		bucket: bucket,
		key:    key,
	}, nil
}

// NewS3ConfigSourceFromEnv creates a new S3 config source using environment variables:
// AWS_REGION, GO_CONFIG_S3_BUCKET, GO_CONFIG_S3_KEY
func NewS3ConfigSourceFromEnv() (ConfigSource, error) {

	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return nil, errors.New("missing AWS_REGION")
	}

	bucket, ok := os.LookupEnv("GO_CONFIG_S3_BUCKET")
	if !ok {
		return nil, errors.New("missing GO_CONFIG_S3_BUCKET")
	}

	key, ok := os.LookupEnv("GO_CONFIG_S3_KEY")
	if !ok {
		return nil, errors.New("missing GO_CONFIG_S3_KEY")
	}

	return NewS3ConfigSource(bucket, key, &region)
}

// Load config file from S3 and pass it to a ViperConfig.
func (source *S3ConfigSource) Load() (Config, error) {

	config := viper.New()
	config.SetConfigType("yaml")

	reader, err := source.readConfig()
	if err != nil {
		return nil, err
	}

	return newViperConfigFromReader(reader)
}

// readConfig downloads the config file from AWS S3 bucket and returns it as an io.Reader.
func (source *S3ConfigSource) readConfig() (io.Reader, error) {

	client := s3.NewFromConfig(source.cfg)
	downloader := manager.NewDownloader(client)

	buf := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(context.TODO(), buf, &s3.GetObjectInput{
		Bucket: aws.String(source.bucket),
		Key:    aws.String(source.key),
	})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
