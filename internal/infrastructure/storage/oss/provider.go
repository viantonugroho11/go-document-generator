// Package oss menyediakan storage provider untuk Alibaba Cloud OSS.
// OSS mendukung S3-compatible API sehingga bisa menggunakan minio-go.
// Compose menggunakan UploadPartCopy (setara S3 multipart) secara server-side.
package oss

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

// NewProvider membuat storage provider untuk Alibaba Cloud OSS.
// endpoint: OSS endpoint S3-compatible, contoh "oss-ap-southeast-5.aliyuncs.com"
// accessKey / secretKey: Alibaba Cloud AccessKeyId dan AccessKeySecret.
// Aktifkan S3-compatible API di Alibaba Console terlebih dahulu.
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
	// Tidak ada enum OSS — pakai S3 karena protokolnya identik.
	// Tambahkan StorageProviderOSS di enums/types.go jika perlu dibedakan di DB.
	return enums.StorageProviderS3
}

func (p *provider) Zip(ctx context.Context, documentID int64, requestID string, entries []sharedStorage.ZipEntry) (string, string, error) {
	return p.inner.Zip(ctx, documentID, requestID, entries)
}

func (p *provider) Compose(ctx context.Context, documentID int64, requestID string, srcPaths []string, ext string) (string, string, error) {
	return p.inner.Compose(ctx, documentID, requestID, srcPaths, ext)
}
