package shared_image_infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mik3lon/starter-template/pkg/file"
)

type S3ImageUploader struct {
	s3Client   *s3.S3
	bucket     string
	s3Endpoint string
	l          Logger
}

func NewS3ImageUploader(c *s3.S3, bucket string, s3Endpoint string, l Logger) *S3ImageUploader {
	return &S3ImageUploader{s3Client: c, bucket: bucket, s3Endpoint: s3Endpoint, l: l}
}

func (s *S3ImageUploader) Upload(ctx context.Context, fi file.FileInfo) (*file.UploadFile, error) {
	_, err := s.createBucketIfNotExists(ctx)
	if err != nil {
		s.l.Warn(ctx, "error creating bucket", map[string]interface{}{
			"error": err.Error(),
		})
	}

	acl := "public-read"
	key := fmt.Sprintf("incidents/%s/images/%s", "test", fi.Filename)
	_, err = s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fi.Content),
		ACL:    &acl,
	})
	if err != nil {
		s.l.Error(ctx, "error uploading image", map[string]interface{}{
			"error": err.Error(),
		})
	}

	url := fmt.Sprintf("%s/%s/%s", s.s3Endpoint, s.bucket, key)

	return &file.UploadFile{
		Url:         url,
		ContentType: fi.ContentType,
		Size:        fi.Size,
		Name:        fi.Filename,
	}, nil
}

func (s *S3ImageUploader) createBucketIfNotExists(ctx context.Context) (*s3.CreateBucketOutput, error) {
	_, err := s.s3Client.HeadBucket(&s3.HeadBucketInput{Bucket: aws.String(s.bucket)})

	if err != nil {
		_, err := s.s3Client.CreateBucket(
			&s3.CreateBucketInput{
				Bucket: aws.String(s.bucket),
			})
		if err != nil {
			s.l.Error(ctx, "error creating bucket", map[string]interface{}{"error": err.Error()})
		}

		err = s.s3Client.WaitUntilBucketExists(&s3.HeadBucketInput{
			Bucket: aws.String(s.bucket),
		})
		if err != nil {
			s.l.Error(ctx, "error occurred while waiting for bucket to be created", map[string]interface{}{"error": err.Error()})
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
			s.l.Error(ctx, "failed to set bucket policy", map[string]interface{}{"error": err.Error()})
		}
	}

	return nil, err
}
