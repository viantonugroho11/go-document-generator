// Package s3 menyediakan storage provider untuk AWS S3.
// Menggunakan minio-go yang mendukung protokol S3 natively.
// ComposeObject di S3 menggunakan CreateMultipartUpload + UploadPartCopy (server-side, tanpa download).
package s3

import (
	"context"
	"time"

	miniostg "go-document-generator/internal/infrastructure/storage/minio"
	"go-document-generator/internal/entity/enums"
	sharedStorage "go-document-generator/internal/shared/storage"
)

type provider struct {
	inner sharedStorage.Provider
}

// NewProvider membuat storage provider untuk AWS S3.
// endpoint: regional endpoint S3, contoh "s3.ap-southeast-1.amazonaws.com"
// accessKey / secretKey: AWS Access Key ID dan Secret Access Key.
// useSSL harus true untuk AWS S3 production.
func NewProvider(endpoint, accessKey, secretKey, bucket string, useSSL bool) (sharedStorage.Provider, error) {
	inner, err := miniostg.NewProvider(endpoint, accessKey, secretKey, bucket, useSSL)
	if err != nil {
		return nil, err
	}
	return &provider{inner: inner}, nil
}

func (p *provider) Save(ctx context.Context, documentID int64, requestID, ext string, data []byte) (string, string, error) {
	return p.inner.Save(ctx, documentID, requestID, ext, data)
}

func (p *provider) Download(ctx context.Context, path string) ([]byte, error) {
	return p.inner.Download(ctx, path)
}

func (p *provider) PresignedURL(ctx context.Context, path string, ttl time.Duration) (string, error) {
	return p.inner.PresignedURL(ctx, path, ttl)
}

func (p *provider) ProviderName() enums.StorageProvider {
	return enums.StorageProviderS3
}

func (p *provider) Zip(ctx context.Context, documentID int64, requestID string, entries []sharedStorage.ZipEntry) (string, string, error) {
	return p.inner.Zip(ctx, documentID, requestID, entries)
}

func (p *provider) Compose(ctx context.Context, documentID int64, requestID string, srcPaths []string, ext string) (string, string, error) {
	return p.inner.Compose(ctx, documentID, requestID, srcPaths, ext)
}
