package config

import (
	"bytes"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

// S3ConfigSource loads a YAML config from a file in a AWS S3 bucket.
type S3ConfigSource struct {

	// Aws config for S3 access.
	config *aws.Config

	// Bucket the config file is located in.
	bucket string

	// Path and file name for a config file.
	key string
}

// NewS3ConfigSource returns a new S3 config source which uses passed config file from given S4 bucket.
// If region is nil it will try to get current aws region from environment var AWS_REGION.
func NewS3ConfigSource(bucket, key string, region *string) ConfigSource {

	if region == nil {
		if envRegion, ok := os.LookupEnv("AWS_REGION"); ok {
			region = &envRegion
		}
	}
	return &S3ConfigSource{
		config: &aws.Config{
			Region: region,
		},
		bucket: bucket,
		key:    key,
	}
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

// readConfig downloads defined config file from AWS S3 bucket and
// uses file content to create a new ViperConfig.
func (source *S3ConfigSource) readConfig() (io.Reader, error) {

	downloader := s3manager.NewDownloader(session.Must(session.NewSession(source.config)))

	buf := &aws.WriteAtBuffer{}
	var reader io.Reader
	_, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(source.bucket),
			Key:    aws.String(source.key),
		})
	if err == nil {
		reader = bytes.NewReader(buf.Bytes())
	}
	return reader, err
}
