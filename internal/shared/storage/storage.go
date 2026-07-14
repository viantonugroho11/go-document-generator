package storage

import (
	"context"
	"time"
)

// Provider abstraksi penyimpanan file dokumen.
type Provider interface {
	// Save menyimpan data ke storage dan mengembalikan path dan nama file.
	Save(ctx context.Context, documentID int64, requestID, ext string, data []byte) (path, fileName string, err error)
	// PresignedURL menghasilkan URL sementara untuk download.
	PresignedURL(ctx context.Context, path string, ttl time.Duration) (string, error)
}
