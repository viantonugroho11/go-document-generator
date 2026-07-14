package storage

import (
	"context"
	"time"

	"go-document-generator/internal/entity/enums"
)

// ZipEntry satu file dalam arsip ZIP.
type ZipEntry struct {
	Name string
	Data []byte
}

// Provider abstraksi penyimpanan file dokumen.
type Provider interface {
	// Save menyimpan data ke storage, mengembalikan path dan nama file.
	Save(ctx context.Context, documentID int64, requestID, ext string, data []byte) (path, fileName string, err error)
	// Download mengambil byte file dari storage berdasarkan path yang tersimpan di DB.
	Download(ctx context.Context, path string) ([]byte, error)
	// PresignedURL menghasilkan URL sementara untuk download.
	PresignedURL(ctx context.Context, path string, ttl time.Duration) (string, error)
	// ProviderName mengembalikan identifier enum provider ini.
	ProviderName() enums.StorageProvider
	// Zip membuat arsip ZIP dari sekumpulan file dan menyimpannya ke storage.
	Zip(ctx context.Context, documentID int64, requestID string, entries []ZipEntry) (path, fileName string, err error)
	// Compose menggabungkan beberapa file dalam storage menjadi satu.
	// Provider yang mendukung server-side compose (MinIO, GCS) melakukannya tanpa download.
	// Provider lokal membaca setiap file dan menggabungkan byte-nya.
	// CATATAN: untuk PDF gunakan library seperti pdfcpu — byte concat tidak menghasilkan PDF valid.
	Compose(ctx context.Context, documentID int64, requestID string, srcPaths []string, ext string) (path, fileName string, err error)
}
