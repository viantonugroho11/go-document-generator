package storage

import (
	"context"
	"fmt"
	"time"
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

func (p *localProvider) PresignedURL(_ context.Context, path string, _ time.Duration) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}
	return path, nil
}
