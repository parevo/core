package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/parevo/core/blob"
)

// Config holds S3 settings.
type Config struct {
	Region          string // AWS region (e.g. us-east-1)
	AccessKeyID     string // AWS access key
	SecretAccessKey string // AWS secret key
}

// Store implements blob.Store for Amazon S3.
type Store struct {
	client *s3.Client
}

// NewStore creates an S3 store.
func NewStore(cfg Config) (*Store, error) {
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("s3: AccessKeyID and SecretAccessKey are required")
	}
	awsCfg := aws.Config{
		Region:      cfg.Region,
		Credentials: credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	}
	return &Store{
		client: s3.NewFromConfig(awsCfg),
	}, nil
}

// Put uploads an object.
func (s *Store) Put(ctx context.Context, bucket, key string, body io.Reader, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	}
	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("s3 put: %w", err)
	}
	return nil
}

// Get downloads an object.
func (s *Store) Get(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 get: %w", err)
	}
	return out.Body, nil
}

// Delete removes an object.
func (s *Store) Delete(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3 delete: %w", err)
	}
	return nil
}

// List returns objects with the given prefix.
func (s *Store) List(ctx context.Context, bucket, prefix string) ([]blob.ObjectInfo, error) {
	out, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("s3 list: %w", err)
	}
	infos := make([]blob.ObjectInfo, 0, len(out.Contents))
	for _, o := range out.Contents {
		info := blob.ObjectInfo{
			Key:  aws.ToString(o.Key),
			Size: aws.ToInt64(o.Size),
		}
		if o.LastModified != nil {
			info.LastModified = *o.LastModified
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// PresignGet returns a presigned URL for downloading the object.
func (s *Store) PresignGet(ctx context.Context, bucket, key string, exp time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	presigned, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(exp))
	if err != nil {
		return "", fmt.Errorf("s3 presign get: %w", err)
	}
	return presigned.URL, nil
}

// PresignPut returns a presigned URL for uploading the object (client uses PUT with the URL).
func (s *Store) PresignPut(ctx context.Context, bucket, key, contentType string, exp time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}
	presigned, err := presigner.PresignPutObject(ctx, input, s3.WithPresignExpires(exp))
	if err != nil {
		return "", fmt.Errorf("s3 presign put: %w", err)
	}
	return presigned.URL, nil
}

var _ blob.PresignedStore = (*Store)(nil)
