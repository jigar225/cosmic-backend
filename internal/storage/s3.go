package storage

import (
	"context"
	"os"
)

// S3Config holds optional S3 configuration from env.
type S3Config struct {
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

// FromEnv reads S3-related env vars. All optional.
func S3FromEnv() S3Config {
	return S3Config{
		Region:    os.Getenv("AWS_REGION"),
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Bucket:    os.Getenv("S3_BUCKET"),
	}
}

// UploadURL returns a placeholder; implement with aws-sdk when needed.
func (c S3Config) UploadURL(ctx context.Context, key string) (string, error) {
	if c.Bucket == "" {
		return "", nil
	}
	_ = key
	return "", nil
}
