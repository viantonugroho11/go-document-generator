package documentrenderlogs

import (
	"context"

	logEntity "go-document-generator/internal/entity/documentrenderlogs"
	docrepo "go-document-generator/internal/repository/documents"
	logrepo "go-document-generator/internal/repository/documentrenderlogs"
	"go-document-generator/internal/shared/pagination"
)

type Service interface {
	ListByDocumentID(ctx context.Context, documentID int64, page pagination.Params) ([]logEntity.RenderLog, pagination.Meta, error)
}

type service struct {
	logs logrepo.DocumentRenderLogsRepository
	docs docrepo.DocumentsRepository
}

func NewService(logs logrepo.DocumentRenderLogsRepository, docs docrepo.DocumentsRepository) Service {
	return &service{logs: logs, docs: docs}
}

func (s *service) ListByDocumentID(ctx context.Context, documentID int64, page pagination.Params) ([]logEntity.RenderLog, pagination.Meta, error) {
	if _, err := s.docs.GetByID(ctx, nil, documentID, nil); err != nil {
		return nil, pagination.Meta{}, err
	}
	page = pagination.Normalize(page.Page, page.Limit)
	items, total, err := s.logs.ListByDocumentID(ctx, nil, documentID, page)
	if err != nil {
		return nil, pagination.Meta{}, err
	}
	return items, pagination.Meta{Page: page.Page, Limit: page.Limit, Total: total}, nil
}
