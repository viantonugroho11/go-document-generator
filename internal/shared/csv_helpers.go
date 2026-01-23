package shared

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CSVQuote melakukan quoting sesuai kebutuhan CSV sederhana:
// - Menggandakan tanda kutip ganda
// - Membungkus dengan tanda kutip bila ada koma, kutip, atau newline
func CSVQuote(s string) string {
	needsQuote := strings.ContainsAny(s, ",\"\n\r")
	if strings.Contains(s, "\"") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
	}
	if needsQuote {
		return `"` + s + `"`
	}
	return s
}

// CSVJoin menggabungkan slice string menjadi satu baris CSV.
// Setiap elemen akan di-quote sesuai aturan CSVQuote.
func CSVJoin(items []string) string {
	quoted := make([]string, len(items))
	for i, it := range items {
		quoted[i] = CSVQuote(it)
	}
	return strings.Join(quoted, ",")
}

// CSVString memformat berbagai tipe umum ke string.
func CSVString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case time.Time:
		return t.Format(time.RFC3339)
	case *time.Time:
		if t == nil {
			return ""
		}
		return t.Format(time.RFC3339)
	case fmt.Stringer:
		return t.String()
	case bool:
		if t {
			return "true"
		}
		return "false"
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", t)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", t)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", t)
	}
}

// DefaultCSVFuncMap menyediakan fungsi-fungsi helper untuk text/template.
func DefaultCSVFuncMap() map[string]any {
	return map[string]any{
		"csvQuote": CSVQuote,
		"csvJoin":  CSVJoin,
		"csvStr":   CSVString,
	}
}
