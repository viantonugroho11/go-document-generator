package minio

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	sharedStorage "go-document-generator/internal/shared/storage"
)

type provider struct {
	client *miniogo.Client
	bucket string
}

// NewProvider membuat storage provider berbasis MinIO/S3-compatible.
func NewProvider(endpoint, accessKey, secretKey, bucket string, useSSL bool) (sharedStorage.Provider, error) {
	client, err := miniogo.New(endpoint, &miniogo.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: init client: %w", err)
	}
	return &provider{client: client, bucket: bucket}, nil
}

func (p *provider) Save(ctx context.Context, documentID int64, requestID, ext string, data []byte) (string, string, error) {
	ext = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(ext)), ".")
	if ext == "" {
		ext = "bin"
	}
	safeReq := sanitize(requestID)
	fileName := fmt.Sprintf("%s.%s", safeReq, ext)
	objectName := fmt.Sprintf("%d/%s", documentID, fileName)

	_, err := p.client.PutObject(ctx, p.bucket, objectName, bytes.NewReader(data), int64(len(data)), miniogo.PutObjectOptions{
		ContentType: contentTypeForExt(ext),
	})
	if err != nil {
		return "", "", fmt.Errorf("minio: put object: %w", err)
	}
	return objectName, fileName, nil
}

func (p *provider) PresignedURL(ctx context.Context, path string, ttl time.Duration) (string, error) {
	u, err := p.client.PresignedGetObject(ctx, p.bucket, path, ttl, nil)
	if err != nil {
		return "", fmt.Errorf("minio: presign: %w", err)
	}
	return u.String(), nil
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "document"
	}
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

func contentTypeForExt(ext string) string {
	switch ext {
	case "pdf":
		return "application/pdf"
	case "html":
		return "text/html"
	case "docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "csv":
		return "text/csv"
	default:
		return "application/octet-stream"
	}
}
