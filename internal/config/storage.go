package config

// Storage konfigurasi provider penyimpanan file dokumen.
type Storage struct {
	// Provider: "local", "minio", "s3", "gcs"
	Provider  string `json:"provider"`
	BaseDir   string `json:"base_dir"`   // untuk local
	Endpoint  string `json:"endpoint"`   // untuk minio/s3
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	UseSSL    bool   `json:"use_ssl"`
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
