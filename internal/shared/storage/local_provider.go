package storage

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"go-document-generator/internal/entity/enums"
)

type localProvider struct {
	baseDir string
}

// NewLocalProvider membuat storage provider berbasis disk lokal.
func NewLocalProvider(baseDir string) Provider {
	if baseDir == "" {
		baseDir = "./storage/documents"
	}
	return &localProvider{baseDir: baseDir}
}

func (p *localProvider) Save(_ context.Context, documentID int64, requestID, ext string, data []byte) (string, string, error) {
	return SaveDocument(p.baseDir, documentID, requestID, ext, data)
}

func (p *localProvider) Download(_ context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}

// PresignedURL untuk local provider mengembalikan path filesystem.
// Download handler harus deteksi ini dan stream file langsung, bukan redirect.
func (p *localProvider) PresignedURL(_ context.Context, path string, _ time.Duration) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}
	return path, nil
}

func (p *localProvider) ProviderName() enums.StorageProvider {
	return enums.StorageProviderLocal
}

func (p *localProvider) Zip(_ context.Context, documentID int64, requestID string, entries []ZipEntry) (string, string, error) {
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
	return SaveDocument(p.baseDir, documentID, requestID, "zip", buf.Bytes())
}

// Compose menggabungkan byte file secara berurutan.
// Hanya cocok untuk format teks (HTML, CSV). Untuk PDF gunakan pdfcpu.
func (p *localProvider) Compose(_ context.Context, documentID int64, requestID string, srcPaths []string, ext string) (string, string, error) {
	var buf bytes.Buffer
	for _, sp := range srcPaths {
		data, err := os.ReadFile(sp)
		if err != nil {
			return "", "", fmt.Errorf("compose: read %s: %w", sp, err)
		}
		buf.Write(data)
	}
	return SaveDocument(p.baseDir, documentID, requestID, ext, buf.Bytes())
}
