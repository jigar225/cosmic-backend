package storage

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Config holds optional S3 configuration from env.
type S3Config struct {
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

// S3FromEnv reads S3-related env vars. All optional.
func S3FromEnv() S3Config {
	return S3Config{
		Region:    os.Getenv("AWS_REGION"),
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Bucket:    os.Getenv("S3_BUCKET"),
	}
}

// S3Uploader uploads files to S3.
type S3Uploader struct {
	client *s3.Client
	bucket string
}

// NewS3Uploader creates an uploader using env credentials (AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, S3_BUCKET).
// If bucket is empty, Upload is a no-op and returns nil (useful for local dev without S3).
func NewS3Uploader(cfg S3Config) (*S3Uploader, error) {
	if cfg.Bucket == "" {
		return &S3Uploader{bucket: ""}, nil
	}
	region := cfg.Region
	if region == "" {
		region = "us-east-1"
	}
	var opts []func(*config.LoadOptions) error
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)))
	}
	awsCfg, err := config.LoadDefaultConfig(context.Background(), append(opts, config.WithRegion(region))...)
	if err != nil {
		return nil, err
	}
	return &S3Uploader{
		client: s3.NewFromConfig(awsCfg),
		bucket: cfg.Bucket,
	}, nil
}

// Upload uploads body to S3 at the given key with contentType and known contentLength.
// Key is stored in DB (e.g. chapters/123/uuid.pdf).
func (u *S3Uploader) Upload(ctx context.Context, key, contentType string, body io.Reader, contentLength int64) error {
	if u.bucket == "" || u.client == nil {
		return nil
	}
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        body,
		ContentLength: aws.Int64(contentLength),
	})
	return err
}

// Bucket returns the configured bucket name (for building URLs if needed).
func (u *S3Uploader) Bucket() string {
	return u.bucket
}

// PresignGetURL returns a presigned URL for downloading the object at key. Valid for 15 minutes.
func (u *S3Uploader) PresignGetURL(ctx context.Context, key string) (string, error) {
	if u.bucket == "" || u.client == nil {
		return "", nil
	}
	presignClient := s3.NewPresignClient(u.client)
	presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", err
	}
	return presignedReq.URL, nil
}
