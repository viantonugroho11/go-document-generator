package pdf

import (
	"bytes"
	"context"
	"html/template"
	"strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// WKHTMLToPDFGenerator mengubah HTML menjadi PDF menggunakan wkhtmltopdf.
// HTML dirender dari template string menggunakan html/template engine.
type WKHTMLToPDFGenerator struct {
	// opsi dasar, dapat ditambah sesuai kebutuhan
	pageSize    string
	orientation string
	dpi         uint
}

// NewWKHTMLToPDFGenerator membuat generator PDF berbasis wkhtmltopdf.
// Gunakan templating.EngineHTML untuk htmlEngine agar auto-escaping HTML aman.
func NewWKHTMLToPDFGenerator() *WKHTMLToPDFGenerator {
	return &WKHTMLToPDFGenerator{
		pageSize:    wkhtml.PageSizeA4,
		orientation: wkhtml.OrientationPortrait,
		dpi:         96,
	}
}

func (g *WKHTMLToPDFGenerator) Generate(ctx context.Context, templateSource string, data any) ([]byte, string, error) {
	// 1) Render HTML dari template (html/template)
	tpl := template.New("html")
	tpl, err := tpl.Parse(templateSource)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, "", err
	}
	html := buf.String()

	// 2) Siapkan wkhtmltopdf
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return nil, "", err
	}
	pdfg.Dpi.Set(g.dpi)
	pdfg.Orientation.Set(g.orientation)
	pdfg.PageSize.Set(g.pageSize)

	page := wkhtml.NewPageReader(strings.NewReader(html))
	pdfg.AddPage(page)

	// 3) Generate PDF
	if err := pdfg.Create(); err != nil {
		return nil, "", err
	}

	return pdfg.Bytes(), "application/pdf", nil
}
