package html

import (
	"bytes"
	"context"
	"html/template"
)

// Generator merender HTML dari template string (engine HTML).
type Generator struct{}

func NewGenerator() *Generator { return &Generator{} }

func (g *Generator) Generate(ctx context.Context, templateSource string, data any) ([]byte, string, error) {
	_ = ctx
	tpl := template.New("html")
	tpl, err := tpl.Parse(templateSource)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "text/html", nil
}
