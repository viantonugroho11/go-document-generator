package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultBaseDir = "./storage/documents"

// SaveDocument menyimpan byte hasil render ke disk lokal.
func SaveDocument(baseDir string, documentID int64, requestID, ext string, data []byte) (absPath string, fileName string, err error) {
	if baseDir == "" {
		baseDir = defaultBaseDir
	}
	ext = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(ext)), ".")
	if ext == "" {
		ext = "bin"
	}
	safeRequest := sanitizeFilePart(requestID)
	fileName = fmt.Sprintf("%s.%s", safeRequest, ext)
	dir := filepath.Join(baseDir, fmt.Sprintf("%d", documentID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", "", err
	}
	absPath = filepath.Join(dir, fileName)
	if err := os.WriteFile(absPath, data, 0o644); err != nil {
		return "", "", err
	}
	return absPath, fileName, nil
}

func sanitizeFilePart(s string) string {
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

func ExtensionForFormat(format string) string {
	switch strings.ToUpper(strings.TrimSpace(format)) {
	case "PDF":
		return "pdf"
	case "HTML":
		return "html"
	case "DOCX":
		return "docx"
	default:
		return "out"
	}
}
