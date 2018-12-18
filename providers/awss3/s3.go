package awss3

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bloom42/rocket/config"
	"github.com/bloom42/astro-go/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/z0mbie42/fswalk"
)

// Deploy perform the S3 upload
func Deploy(conf config.AWSS3Config) error {
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

	if conf.LocalDirectory == nil {
		v := "."
		conf.LocalDirectory = &v
	}

	if conf.RemoteDirectory == nil {
		v := "/"
		conf.RemoteDirectory = &v
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

	walker, _ := fswalk.NewWalker()
	filesc, _ := walker.Walk(*conf.LocalDirectory)
	for file := range filesc {
		if file.Path == "." || file.IsDir || file.IsSymLink {
			continue
		}
		log.With("file", file.Path).Debug("aws_s3: file to upload")
		err = UploadFileToS3(conf, sess, file.Path)
		if err != nil {
			log.With("file", file.Path).Error(fmt.Sprintf("aws_s3: error uploading a file: %s", err.Error()))
		} else {
			log.Info(fmt.Sprintf("aws_s3: file successfully uploaded %s", file.Path))
		}
	}
	return nil
}

func UploadFileToS3(conf config.AWSS3Config, s *session.Session, filePath string) error {

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
