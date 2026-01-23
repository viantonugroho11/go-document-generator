package csv

import (
	"bytes"
	"context"
	"text/template"

	"go-document-generator/internal/shared"
)

// TmplCSVGenerator merender CSV menggunakan text/template (engine "tmpl").
// Gunakan helper DefaultCSVFuncMap untuk utilitas umum CSV.
type TmplCSVGenerator struct {
	funcs map[string]any
}

// NewCSVGenerator membuat generator CSV berbasis text/template.
// Opsi funcs akan digabung dengan DefaultCSVFuncMap.
func NewCSVGenerator() *TmplCSVGenerator {
	funcs := shared.DefaultCSVFuncMap()
	return &TmplCSVGenerator{
		funcs: funcs,
	}
}

func (g *TmplCSVGenerator) Generate(ctx context.Context, templateSource string, data any) ([]byte, string, error) {
	tpl := template.New("csv").Funcs(template.FuncMap(g.funcs))
	tpl, err := tpl.Parse(templateSource)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "text/csv", nil
}
