# Storage Providers

Dokumentasi tentang cara kerja penyimpanan file dokumen, termasuk operasi Zip dan Merge.

---

## Arsitektur

Storage menggunakan interface abstrak `shared/storage.Provider` sehingga seluruh usecase layer tidak bergantung pada implementasi konkret.

```
usecase/documents.StorageProvider  ← subset interface (PresignedURL, Download, Zip, Compose)
         │
shared/storage.Provider            ← interface lengkap (+ Save, ProviderName)
         │
         ├── localProvider          ← disk lokal (dev)
         ├── minio.provider         ← MinIO / S3-compatible (prod)
         └── gcs.provider           ← Google Cloud Storage (belum diimplementasi, lihat bawah)
```

---

## Interface `storage.Provider`

```go
type Provider interface {
    Save(ctx, documentID, requestID, ext, data)  (path, fileName, error)
    Download(ctx, path)                           ([]byte, error)
    PresignedURL(ctx, path, ttl)                  (string, error)
    ProviderName()                                enums.StorageProvider
    Zip(ctx, documentID, requestID, entries)      (path, fileName, error)
    Compose(ctx, documentID, requestID, srcPaths, ext) (path, fileName, error)
}
```

---

## Local Provider

**Pakai**: development, testing.

| Metode       | Perilaku |
|--------------|----------|
| `Save`       | Tulis ke `{baseDir}/{documentID}/{requestID}.{ext}` |
| `Download`   | `os.ReadFile(path)` |
| `PresignedURL` | Kembalikan path filesystem apa adanya |
| `Zip`        | Buat ZIP dalam memori (`archive/zip`), simpan ke disk |
| `Compose`    | Baca setiap file, gabungkan byte, simpan — **hanya cocok untuk HTML/CSV** |

> **Catatan Download handler**: jika URL tidak dimulai `http://`/`https://`, handler langsung stream file via `c.File(path)` sehingga client tidak perlu akses filesystem server.

**Config** (`configs/config.yaml`):
```yaml
storage:
  provider: local
  base_dir: ./storage/documents
```

---

## MinIO Provider

**Pakai**: production, staging, atau lokal dengan `docker compose`.

| Metode       | Perilaku |
|--------------|----------|
| `Save`       | `PutObject` ke bucket di `{documentID}/{fileName}` |
| `Download`   | `GetObject` → baca semua byte |
| `PresignedURL` | `PresignedGetObject` dengan TTL |
| `Zip`        | Buat ZIP dalam memori, `PutObject` sekali |
| `Compose`    | `ComposeObject` (server-side, tanpa download tiap chunk) |

**Config**:
```yaml
storage:
  provider: minio
  endpoint: localhost:9000
  access_key: minioadmin
  secret_key: minioadmin
  bucket: documents
  use_ssl: false
```

### Keterbatasan `Compose` untuk PDF
MinIO `ComposeObject` menggabungkan byte object secara literal — ini **bukan** PDF merge yang valid. PDF mempunyai struktur header/footer/xref yang harus direkonstruksi. Gunakan `pdfcpu` untuk merge PDF yang benar:

```go
import "github.com/pdfcpu/pdfcpu/pkg/api"

// Merge PDF dari file paths
err := api.MergeCreateFile(srcPaths, destPath, conf)
```

---

## GCS Provider (belum diimplementasi)

Google Cloud Storage menyediakan `Compose` API yang setara dengan MinIO server-side compose.

### Cara implementasi

```go
// internal/infrastructure/storage/gcs/provider.go
package gcs

import (
    "cloud.google.com/go/storage"
    sharedStorage "go-document-generator/internal/shared/storage"
    "go-document-generator/internal/entity/enums"
)

type provider struct {
    client *storage.Client
    bucket string
}

func (p *provider) Compose(ctx context.Context, documentID int64, requestID string, srcPaths []string, ext string) (string, string, error) {
    bkt := p.client.Bucket(p.bucket)
    dstName := fmt.Sprintf("%d/%s.%s", documentID, requestID, ext)
    dst := bkt.Object(dstName)

    srcs := make([]*storage.ObjectHandle, len(srcPaths))
    for i, sp := range srcPaths {
        srcs[i] = bkt.Object(sp)
    }
    composer := dst.ComposerFrom(srcs...)
    if _, err := composer.Run(ctx); err != nil {
        return "", "", fmt.Errorf("gcs: compose: %w", err)
    }
    return dstName, filepath.Base(dstName), nil
}

func (p *provider) ProviderName() enums.StorageProvider {
    return enums.StorageProviderGCS // tambahkan enum ini
}
```

**Config**:
```yaml
storage:
  provider: gcs
  bucket: my-documents-bucket
  # credentials via GOOGLE_APPLICATION_CREDENTIALS env var
```

> GCS `Compose` dibatasi **32 object per operasi**. Untuk lebih dari 32 file, lakukan compose bertahap (composing compositions).

---

## Operasi Zip

### Endpoint

```
POST /documents/zip
Authorization: X-API-Key: <key>
Content-Type: application/json

{
  "ids": [1, 2, 3],
  "label": "laporan-q1"   // opsional, dipakai sebagai nama file
}
```

### Response

- **MinIO / GCS**: `200 OK` dengan `{"url": "https://...?X-Amz-..."}` (presigned URL, valid 15 menit)
- **Local**: file langsung distream ke client

### Flow

```
handler → service.ZipDocuments(ids)
           ├── docs.GetByID() per ID (validasi status = GENERATED)
           ├── storage.Download(filePath) per dokumen
           ├── storage.Zip(entries)      → simpan ke storage
           └── storage.PresignedURL()   → URL download
```

---

## Operasi Merge

### Endpoint

```
POST /documents/merge
Authorization: X-API-Key: <key>
Content-Type: application/json

{
  "ids": [1, 2],
  "label": "merged-contract"   // opsional
}
```

**Constraint**: semua dokumen harus format yang sama (misalnya semua HTML atau semua PDF).

### Response

Sama dengan Zip — URL atau stream tergantung provider.

### Flow

```
handler → service.MergeDocuments(ids)
           ├── docs.GetByID() per ID
           ├── validasi semua format sama
           ├── storage.Compose(srcPaths, ext)
           │     ├── MinIO: ComposeObject (server-side, efisien)
           │     ├── GCS:   ComposerFrom (server-side, max 32 object)
           │     └── Local: os.ReadFile + byte concat
           └── storage.PresignedURL()
```

### Merge PDF yang benar

Saat ini `Compose` untuk local dan MinIO melakukan byte concat, yang **tidak menghasilkan PDF valid**. Untuk PDF:

1. Tambahkan dependency:
   ```
   go get github.com/pdfcpu/pdfcpu
   ```

2. Override implementasi `Compose` untuk ext `pdf`:
   ```go
   if ext == "pdf" {
       return composePDF(ctx, documentID, requestID, srcPaths)
   }
   ```

---

## Menambah Provider Baru

1. Buat package di `internal/infrastructure/storage/<provider>/provider.go`
2. Implementasikan semua method `shared/storage.Provider`
3. Tambahkan enum di `internal/entity/enums/storage_provider.go`
4. Daftarkan di `internal/bootstrap/wire_documents.go`:
   ```go
   case "gcs":
       sp, spErr = gcsstg.NewProvider(c.Storage.Bucket, c.Storage.CredentialsFile)
   ```
5. Tambahkan config key di `configs/config.yaml` dan `internal/config/storage.go`
