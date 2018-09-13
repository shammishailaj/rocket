package eb

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/astrocorp42/rocket/config"
	"github.com/astroflow/astroflow-go/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
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

	// 1) create the archive
	tmpFile, err := ioutil.TempFile("", "rocket.*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	// set up the zip writer
	zipWriter := zip.NewWriter(tmpFile)

	walker, _ := fswalk.NewWalker()
	filesc, _ := walker.Walk(*conf.Directory)
	for file := range filesc {
		if file.Path == "." || file.IsDir || file.IsSymLink {
			continue
		}
		log.With("bundle", tmpFile.Name(), "file", file.Path).Debug("aws_eb: adding file to bundle")
		err = addFileToBundle(zipWriter, file.Path)
		if err != nil {
			return err
		}
	}

	if err = zipWriter.Close(); err != nil {
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}

	// 2) upload it to S3
	err = uploadFileToS3(conf, sess, tmpFile.Name())
	if err != nil {
		log.With("file", tmpFile.Name()).Error(fmt.Sprintf("aws_eb: error uploading a file to S3: %s", err.Error()))
	} else {
		log.Info(fmt.Sprintf("aws_eb: file successfully uploaded to S3%s", tmpFile.Name()))
	}

	// 3) create a new EB application version
	svc := elasticbeanstalk.New(sess)
	input := &elasticbeanstalk.CreateApplicationVersionInput{
		ApplicationName:       aws.String(*conf.Application),
		AutoCreateApplication: aws.Bool(true),
		Process:               aws.Bool(true),
		SourceBundle: &elasticbeanstalk.S3Location{
			S3Bucket: aws.String(*conf.S3Bucket),
			S3Key:    aws.String(filepath.Join(*conf.S3Directory, filepath.Base(tmpFile.Name()))),
		},
		VersionLabel: aws.String(*conf.Version),
	}

	_, err = svc.CreateApplicationVersion(input)
	return err
}

func addFileToBundle(zw *zip.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(stat)
	if err != nil {
		return err
	}

	header.Name = path
	header.Method = zip.Deflate

	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	if _, err = io.Copy(writer, file); err != nil {
		return err
	}
	return nil
}

func uploadFileToS3(conf config.AWSEBConfig, s *session.Session, filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(*conf.S3Bucket),
		Key:    aws.String(filepath.Join(*conf.S3Directory, filepath.Base(filePath))),
		Body:   file,
	})
	return err
}
