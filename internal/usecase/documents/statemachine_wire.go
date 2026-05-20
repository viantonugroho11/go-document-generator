package documents

import (
	"go-document-generator/internal/usecase/documents/states"
	"go-document-generator/internal/usecase/documents/transitions"
)

// BuildStateHandlers menyusun handler transisi (wiring; hindari states → transitions).
func BuildStateHandlers(deps transitions.Deps) states.Handlers {
	return states.Handlers{
		OnFieldUpdate:  transitions.NewFieldUpdate(deps),
		OnToQueued:     transitions.NewToQueued(),
		OnToProcessing: transitions.NewToProcessing(),
		OnToGenerated:  transitions.NewToGenerated(deps),
		OnToCancelled:  transitions.NewToCancelled(),
		OnToFailed:     transitions.NewToFailed(),
		OnRetry:        transitions.NewRetry(),
		OnTerminal:     transitions.NewNoop(),
	}
}
