package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"time"
)

type bucket string

func (b bucket) string() string {
	return string(b)
}

func (s *service) getPresignURL(bucket bucket, object string, expired time.Duration) (string, error) {
	presignResult, err := s.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket.string()),
		Key:    aws.String(object),
	}, func(options *s3.PresignOptions) {
		options.Expires = expired
	})
	if err != nil {
		return "", err
	}
	return presignResult.URL, err
}

func (s *service) UploadFile(file io.Reader, path string) error {
	_, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.publicBucket.string()),
		Key:    aws.String(path),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to %s/%s: %w", s.publicBucket, path, err)
	}
	return nil
}

func (s *service) CopyFileToAnotherBucket(path string, oldBucket bucket, newBucket bucket) error {
	_, err := s.s3Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(newBucket.string()),
		CopySource: aws.String(fmt.Sprintf("%s/%s", oldBucket.string(), path)),
		Key:        aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to copy object for soft delete: %w", err)
	}
	return nil
}

func (s *service) DeleteFile(path string, b bucket) error {
	_, err := s.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.string()),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s/%s: %w", s.publicBucket, path, err)
	}
	return nil
}

func (s *service) ReadFile(path string) ([]byte, error) {
	output, err := s.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.publicBucket.string()),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s/%s: %w", s.publicBucket, path, err)
	}
	defer output.Body.Close()

	content, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}
	return content, nil
}

func (s *service) ListFiles() error {
	output, err := s.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:                   aws.String(s.publicBucket.string()),
		ContinuationToken:        nil,
		Delimiter:                nil,
		EncodingType:             "",
		MaxKeys:                  nil,
		OptionalObjectAttributes: nil,
		Prefix:                   nil,
		RequestPayer:             "",
		StartAfter:               nil,
	})
	if err != nil {
		return fmt.Errorf("failed to list objects in bucket %s: %w", s.publicBucket, err)
	}

	fmt.Println(output)

	return nil
}
