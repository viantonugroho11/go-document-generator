package documents

import (
	"errors"

	"go-document-generator/internal/entity/enums"
)

// Urutan alur utama: PENDING → QUEUED → PROCESSING → GENERATED.
var statusOrder = []enums.DocumentStatus{
	enums.DocumentStatusPending,
	enums.DocumentStatusQueued,
	enums.DocumentStatusProcessing,
	enums.DocumentStatusGenerated,
}

func statusIndex(s enums.DocumentStatus) int {
	for i, st := range statusOrder {
		if st == s {
			return i
		}
	}
	return -1
}

// canTransitionStatus memvalidasi perpindahan status (harus berurutan atau cabang yang diizinkan).
func canTransitionStatus(from, to enums.DocumentStatus) bool {
	if from == to {
		return true
	}
	if to == enums.DocumentStatusCancelled {
		switch from {
		case enums.DocumentStatusPending, enums.DocumentStatusQueued, enums.DocumentStatusProcessing:
			return true
		default:
			return false
		}
	}
	if to == enums.DocumentStatusFailed {
		return from == enums.DocumentStatusProcessing
	}
	if from == enums.DocumentStatusFailed && to == enums.DocumentStatusQueued {
		return true // retry via patch (alternatif endpoint Retry)
	}

	fromIdx := statusIndex(from)
	toIdx := statusIndex(to)
	if fromIdx < 0 || toIdx < 0 {
		return false
	}
	// Hanya boleh ke status berikutnya dalam urutan (selisih 1).
	return toIdx == fromIdx+1
}

func validateStatusTransition(from, to enums.DocumentStatus) error {
	if canTransitionStatus(from, to) {
		return nil
	}
	return errors.New("invalid status transition: must follow PENDING → QUEUED → PROCESSING → GENERATED")
}

func allowsFieldPatch(status enums.DocumentStatus) bool {
	switch status {
	case enums.DocumentStatusPending, enums.DocumentStatusQueued:
		return true
	default:
		return false
	}
}
