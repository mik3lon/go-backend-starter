package shared_image_infrastructure

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mik3lon/starter-template/pkg/file"
	"log"
)

type S3ImageUploader struct {
	s3Client   *s3.S3
	bucket     string
	s3Endpoint string
}

func NewS3ImageUploader(c *s3.S3, bucket string, s3Endpoint string) *S3ImageUploader {
	return &S3ImageUploader{s3Client: c, bucket: bucket, s3Endpoint: s3Endpoint}
}

func (s *S3ImageUploader) Upload(fi file.FileInfo) (*file.UploadFile, error) {
	s.createBucketIfNotExists()

	acl := "public-read"
	key := fmt.Sprintf("incidents/%s/images/%s", "test", fi.Filename)
	_, err := s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fi.Content),
		ACL:    &acl,
	})
	if err != nil {
		log.Fatalf("Unable to upload %q to %q, %v", key, s.bucket, err)
	}

	url := fmt.Sprintf("%s/%s/%s", s.s3Endpoint, s.bucket, key)

	return &file.UploadFile{
		Url:         url,
		ContentType: fi.ContentType,
		Size:        fi.Size,
		Name:        fi.Filename,
	}, nil
}

func (s *S3ImageUploader) createBucketIfNotExists() (*s3.CreateBucketOutput, error) {
	_, err := s.s3Client.HeadBucket(
		&s3.HeadBucketInput{Bucket: aws.String(s.bucket)})

	if err != nil {
		return s.s3Client.CreateBucket(
			&s3.CreateBucketInput{
				Bucket: aws.String(s.bucket),
			})
	}

	err = s.s3Client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		log.Fatalf("Error occurred while waiting for bucket to be created: %v", err)
	}

	bucketPolicy := `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "PublicReadGetObject",
                "Effect": "Allow",
                "Principal": "*",
                "Action": "s3:GetObject",
                "Resource": "arn:aws:s3:::%s/*"
            }
        ]
    }`
	policy := fmt.Sprintf(bucketPolicy, s.bucket)
	_, err = s.s3Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(s.bucket),
		Policy: aws.String(policy),
	})
	if err != nil {
		log.Fatalf("Failed to set bucket policy: %v", err)
	}

	return nil, err
}
