package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"keeper/internal/core/config"
	"keeper/internal/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type S3Client struct {
	s3client   *minio.Client
	bucketName string
}

func NewS3Client(ctx context.Context, cfg *config.S3Config) (*S3Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start new s3 server: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket (%s) exists: %w", cfg.BucketName, err)
	}
	if !exists {
		logger.Log.Debug("Bucket not found, create backet.", zap.String("BacketName", cfg.BucketName))
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create new bucket: %w", err)
		}
		logger.Log.Debug("Successfully created", zap.String("BacketName", cfg.BucketName))
	}

	return &S3Client{s3client: client, bucketName: cfg.BucketName}, nil
}

func (s *S3Client) PutObject(ctx context.Context, name string, obj io.Reader, size int64) error {
	if _, err := s.s3client.PutObject(ctx, s.bucketName, name, obj, size, minio.PutObjectOptions{}); err != nil {
		return fmt.Errorf("failed to put object into bucket: %w", err)
	}

	return nil
}

func (s *S3Client) GetObject(ctx context.Context, name string) ([]byte, error) {
	obj, err := s.s3client.GetObject(ctx, s.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	defer func() {
		if err = obj.Close(); err != nil {
			logger.Log.Error("failed to close object", zap.String("objectName", name))
		}
	}()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(obj); err != nil {
		return nil, fmt.Errorf("failed to read obj from S3: %w", err)
	}
	return buf.Bytes(), nil
}
