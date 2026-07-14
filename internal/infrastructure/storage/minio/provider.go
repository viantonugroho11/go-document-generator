package minio

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go-document-generator/internal/entity/enums"
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

func (p *provider) Download(ctx context.Context, path string) ([]byte, error) {
	obj, err := p.client.GetObject(ctx, p.bucket, path, miniogo.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("minio: get object: %w", err)
	}
	defer obj.Close()
	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("minio: read object: %w", err)
	}
	return data, nil
}

func (p *provider) PresignedURL(ctx context.Context, path string, ttl time.Duration) (string, error) {
	u, err := p.client.PresignedGetObject(ctx, p.bucket, path, ttl, nil)
	if err != nil {
		return "", fmt.Errorf("minio: presign: %w", err)
	}
	return u.String(), nil
}

func (p *provider) ProviderName() enums.StorageProvider {
	return enums.StorageProviderMinio
}

// Zip membuat arsip ZIP dalam memori dari entries, lalu menyimpan ke MinIO.
func (p *provider) Zip(ctx context.Context, documentID int64, requestID string, entries []sharedStorage.ZipEntry) (string, string, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, e := range entries {
		f, err := w.Create(e.Name)
		if err != nil {
			return "", "", fmt.Errorf("zip: create entry %s: %w", e.Name, err)
		}
		if _, err := f.Write(e.Data); err != nil {
			return "", "", fmt.Errorf("zip: write entry %s: %w", e.Name, err)
		}
	}
	if err := w.Close(); err != nil {
		return "", "", fmt.Errorf("zip: close writer: %w", err)
	}
	return p.Save(ctx, documentID, requestID, "zip", buf.Bytes())
}

// Compose menggunakan ComposeObject MinIO (server-side, tanpa download setiap file).
// Cocok untuk semua format. Untuk GCS gunakan storage.Compose API yang setara.
func (p *provider) Compose(ctx context.Context, documentID int64, requestID string, srcPaths []string, ext string) (string, string, error) {
	ext = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(ext)), ".")
	if ext == "" {
		ext = "bin"
	}
	safeReq := sanitize(requestID)
	fileName := fmt.Sprintf("%s.%s", safeReq, ext)
	objectName := fmt.Sprintf("%d/%s", documentID, fileName)

	srcs := make([]miniogo.CopySrcOptions, len(srcPaths))
	for i, sp := range srcPaths {
		srcs[i] = miniogo.CopySrcOptions{Bucket: p.bucket, Object: sp}
	}
	dst := miniogo.CopyDestOptions{Bucket: p.bucket, Object: objectName}
	if _, err := p.client.ComposeObject(ctx, dst, srcs...); err != nil {
		return "", "", fmt.Errorf("minio: compose: %w", err)
	}
	return objectName, fileName, nil
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
	case "zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}
