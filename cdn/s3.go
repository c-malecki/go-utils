package cdn

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/c-malecki/go-utils/gen"
	"github.com/c-malecki/go-utils/img/avatar"
)

/*
	TODO:
	"github.com/aws/aws-sdk-go/aws/session" is deprecated: Deprecated: aws-sdk-go is deprecated. Use aws-sdk-go-v2.
	See https://aws.amazon.com/blogs/developer/announcing-end-of-support-for-aws-sdk-for-go-v1-on-july-31-2025/.
*/

type S3ClientConfig struct {
	PublicURL    string // ex: https://cdn.example.com/ or https://cdn.staging.example.com/
	PublicPrefix string // ex: staging ->  https://cdn.staging.example.com/staging/
	S3Endpoint   string // ex: "https://nyc3.digitaloceanspaces.com"
	S3Region     string // ex: "nyc3"
	S3Key        string
	S3Secret     string
	S3Bucket     string
}

type S3Client struct {
	s3     *s3.S3
	bucket string
	url    string
	prefix string
}

type Cdn interface {
	SetPicture(pictureFilename *string) *string
	UploadFile(fileData []byte, extension string, contentType string) (string, error)
	DeleteFile(objectKey string) error
}

func validateS3Key(key string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_/-]+$`)
	return len(key) >= 1 && len(key) <= 1024 && re.MatchString(key)
}

func CreateS3Client(config S3ClientConfig) (*S3Client, error) {
	if len(config.PublicPrefix) > 0 && !validateS3Key(config.PublicPrefix) {
		return nil, fmt.Errorf("public prefix \"%s\" is invalid: must contain characters \"a-z A-Z 0-9 _ -\" only", config.PublicPrefix)
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.S3Key, config.S3Secret, ""),
		Endpoint:         aws.String(config.S3Endpoint),
		Region:           aws.String(config.S3Region),
		S3ForcePathStyle: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	c := s3.New(sess)

	client := &S3Client{
		s3:     c,
		bucket: config.S3Bucket,
		url:    config.PublicURL,
		prefix: config.PublicPrefix,
	}

	return client, nil
}

func (c *S3Client) SetImage(imageFile string, name string) string {
	if strings.HasPrefix(imageFile, "data:image/svg+xml") {
		return imageFile
	}
	if len(imageFile) > 0 {
		return c.url + imageFile
	}
	return avatar.SVGWithInitials(name)
}

// TODO: configure ACL
func (c *S3Client) UploadFile(data []byte, ext string, contentType string) (string, error) {
	filename := gen.GenerateUniqueFilename(ext)

	objectKey := filename
	if len(c.prefix) > 0 {
		objectKey = objectKey + "/" + filename
	}

	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	_, err := c.s3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(data),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	return filename, nil
}

func (c *S3Client) DeleteFile(filename string) error {
	objectKey := filename
	if len(c.prefix) > 0 {
		objectKey = objectKey + "/" + filename
	}

	_, err := c.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}

	return nil
}
