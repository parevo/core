package r2

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/parevo/core/blob"
)

// Config holds Cloudflare R2 settings.
// R2 is S3-compatible. Create API tokens at: Cloudflare Dashboard > R2 > Manage R2 API Tokens
type Config struct {
	AccountID       string // Cloudflare account ID
	AccessKeyID     string // R2 API token access key
	SecretAccessKey string // R2 API token secret key
}

// Store implements blob.Store for Cloudflare R2 (S3-compatible).
type Store struct {
	client *s3.Client
}

// NewStore creates an R2 store.
func NewStore(cfg Config) (*Store, error) {
	if cfg.AccountID == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("r2: AccountID, AccessKeyID, and SecretAccessKey are required")
	}
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("r2 config: %w", err)
	}
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
	return &Store{client: client}, nil
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
		return fmt.Errorf("r2 put: %w", err)
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
		return nil, fmt.Errorf("r2 get: %w", err)
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
		return fmt.Errorf("r2 delete: %w", err)
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
		return nil, fmt.Errorf("r2 list: %w", err)
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
		return "", fmt.Errorf("r2 presign get: %w", err)
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
		return "", fmt.Errorf("r2 presign put: %w", err)
	}
	return presigned.URL, nil
}

var _ blob.PresignedStore = (*Store)(nil)
