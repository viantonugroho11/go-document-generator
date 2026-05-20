package documents

import "context"

// Generator merender dokumen dari template + data payload.
type Generator interface {
	Generate(ctx context.Context, templateSource string, data any) ([]byte, string, error)
}

// GeneratorSelector memilih engine render berdasarkan format output dan template engine.
type GeneratorSelector interface {
	Select(outputFormat string, engine string) Generator
}
