package documents

import "go-document-generator/internal/usecase/documents/transitions"

type generatorSelectorAdapter struct {
	inner GeneratorSelector
}

func (a generatorSelectorAdapter) Select(outputFormat, engine string) transitions.Generator {
	return a.inner.Select(outputFormat, engine)
}

func adaptSelector(s GeneratorSelector) transitions.GeneratorSelector {
	if s == nil {
		return nil
	}
	return generatorSelectorAdapter{inner: s}
}
