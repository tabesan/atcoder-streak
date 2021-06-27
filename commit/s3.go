package commit

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

const bucketName = "streak-bucket"
const objectKey = "streak.txt"

func (c *Client) NewSession() (*session.Session, error) {
	accessKey := os.Getenv("ACCESS")
	privateAccessKey := os.Getenv("PRI_ACCESS")
	cr := credentials.NewStaticCredentials(accessKey, privateAccessKey, "")
	session, err := session.NewSession(&aws.Config{
		Credentials: cr,
		Region:      aws.String("ap-northeast-1"),
	})
	if err != nil {
		errors.Wrap(err, "Create new session error")
	}

	return session, err
}

func (c *Client) DownloadData() (int, string, string, error) {
	session, err := c.NewSession()
	svc := s3.New(session)
	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return 0, "", "", errors.Wrap(err, "GetObject from aws error")
	}
	rc := obj.Body
	defer rc.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	data := strings.Split(buf.String(), " ")
	streak, err := strconv.Atoi(data[0])
	latest := data[1]
	update := data[2]

	return streak, latest, update, nil
}

func (c *Client) UploadData() error {
	session, err := c.NewSession()
	uploader := s3manager.NewUploader(session)
	newInfo := strings.Join([]string{strconv.Itoa(c.streak), c.latestCommit, c.edit.Today()}, " ")
	err = ioutil.WriteFile("streak.txt", []byte(newInfo), 0664)
	if err != nil {
		return errors.Wrap(err, "WriteFile error")
	}

	file, err := os.Open(objectKey)
	if err != nil {
		errors.Wrap(err, "os.Open error")
		return err
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("streak.txt"),
		Body:   file,
	})
	if err != nil {
		return errors.Wrap(err, "UploadData error")
	}

	return nil
}
