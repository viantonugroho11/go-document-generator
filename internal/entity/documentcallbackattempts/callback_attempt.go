package documentcallbackattempts

import (
	"time"
)

type CallbackAttempt struct {
	ID                 int64
	DocumentID         int64
	CallbackURL        string
	RequestPayload     map[string]any
	ResponsePayload    map[string]any
	ResponseStatusCode *int
	IsSuccess          bool
	ErrorMessage       *string
	AttemptedAt        time.Time
}
