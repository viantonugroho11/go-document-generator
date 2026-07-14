package config

// Storage konfigurasi provider penyimpanan file dokumen.
// Provider yang didukung: "local", "minio", "s3" / "aws", "oss" / "alibaba", "gcs"
type Storage struct {
	Provider  string `json:"provider"`
	BaseDir   string `json:"base_dir"`   // local: direktori dokumen
	Endpoint  string `json:"endpoint"`   // cloud: hostname tanpa protokol
	Region    string `json:"region"`     // aws s3: region, contoh "ap-southeast-1"
	AccessKey string `json:"access_key"` // aws: Access Key ID; alibaba: AccessKeyId
	SecretKey string `json:"secret_key"` // aws: Secret Access Key; alibaba: AccessKeySecret
	Bucket    string `json:"bucket"`
	UseSSL    bool   `json:"use_ssl"`    // true untuk production cloud
}

// Auth konfigurasi autentikasi API.
type Auth struct {
	// APIKeys daftar API key yang valid. Kosong = auth dinonaktifkan (dev mode).
	APIKeys []string `json:"api_keys"`
}

// Dms konfigurasi Document Management System.
type Dms struct {
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key"`
}

// CallbackConfig konfigurasi webhook callback.
type CallbackConfig struct {
	// HMACSecret digunakan untuk sign callback payload. Kosong = tanpa signature.
	HMACSecret string `json:"hmac_secret"`
	MaxRetries int    `json:"max_retries"`
}
