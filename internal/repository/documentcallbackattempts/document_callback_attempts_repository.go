package documentcallbackattempts

import (
	"context"

	cbEntity "go-document-generator/internal/entity/documentcallbackattempts"
	"go-document-generator/internal/shared/pagination"

	"gorm.io/gorm"
)

type DocumentCallbackAttemptsRepository interface {
	Create(ctx context.Context, tx *gorm.DB, a cbEntity.CallbackAttempt) (cbEntity.CallbackAttempt, error)
	ListByDocumentID(ctx context.Context, tx *gorm.DB, documentID int64, page pagination.Params) ([]cbEntity.CallbackAttempt, int64, error)
}
