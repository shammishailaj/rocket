package eb

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astrocorp42/rocket/config"
	"github.com/astroflow/astroflow-go/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/z0mbie42/fswalk"
)

// Deploy perform the elastic beanstalk deployment
func Deploy(conf config.AWSEBConfig) error {
	var err error

	if conf.AccessKeyID == nil {
		v := os.Getenv("AWS_ACCESS_KEY_ID")
		conf.AccessKeyID = &v
	} else {
		v := config.ExpandEnv(*conf.AccessKeyID)
		conf.AccessKeyID = &v
	}

	if conf.SecretAccessKey == nil {
		v := os.Getenv("AWS_SECRET_ACCESS_KEY")
		conf.SecretAccessKey = &v
	} else {
		v := config.ExpandEnv(*conf.SecretAccessKey)
		conf.SecretAccessKey = &v
	}

	if conf.Region == nil {
		v := os.Getenv("AWS_REGION")
		conf.Region = &v
	} else {
		v := config.ExpandEnv(*conf.Region)
		conf.Region = &v
	}

	if conf.Application == nil {
		v := os.Getenv("AWS_EB_APPLICATION")
		conf.Application = &v
	} else {
		v := config.ExpandEnv(*conf.Application)
		conf.Application = &v
	}

	if conf.Environment == nil {
		v := os.Getenv("AWS_EB_ENVIRONMENT")
		conf.Environment = &v
	} else {
		v := config.ExpandEnv(*conf.Environment)
		conf.Environment = &v
	}

	if conf.S3Bucket == nil {
		v := os.Getenv("AWS_S3_BUCKET")
		conf.S3Bucket = &v
	} else {
		v := config.ExpandEnv(*conf.S3Bucket)
		conf.S3Bucket = &v
	}

	if conf.Version == nil {
		v := os.Getenv("ROCKET_COMMIT_HASH")
		conf.Version = &v
	} else {
		v := config.ExpandEnv(*conf.Version)
		conf.Version = &v
	}

	if conf.Directory == nil {
		v := "."
		conf.Directory = &v
	}

	if conf.S3Directory == nil {
		v := "/"
		conf.S3Directory = &v
	}

	var awsConf aws.Config

	if *conf.AccessKeyID != "" && *conf.SecretAccessKey != "" {
		awsConf = aws.Config{
			Credentials: credentials.NewStaticCredentials(*conf.AccessKeyID, *conf.SecretAccessKey, ""),
		}
	} else {
		awsConf = aws.Config{}
	}
	awsConf.Region = aws.String(*conf.Region)
	sess := session.New(&awsConf)

	return nil
}

func uploadFileToS3(conf config.AWSS3Config, s *session.Session, filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(*conf.Bucket),
		Key:    aws.String(filepath.Join(*conf.RemoteDirectory, filepath.Base(filePath))),
		Body:   file,
	})
	return err
}
