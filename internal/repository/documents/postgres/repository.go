package postgres

import (
	"context"
	"errors"
	"time"

	entity "go-document-generator/internal/entity/documents"
	repo "go-document-generator/internal/repository/documents"
	"go-document-generator/internal/repository/documents/model"

	"gorm.io/gorm"
)

type documentsRepository struct {
	db *gorm.DB
}

func NewDocumentsRepository(db *gorm.DB) repo.DocumentsRepository {
	return &documentsRepository{db: db}
}

func (r *documentsRepository) Create(ctx context.Context, doc entity.Document) (entity.Document, error) {
	m := toModel(doc)
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = m.CreatedAt
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return entity.Document{}, err
	}
	return toEntity(m), nil
}

func (r *documentsRepository) GetByID(ctx context.Context, id int64) (entity.Document, error) {
	var m model.Document
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Document{}, errors.New("document not found")
	}
	return toEntity(m), err
}

func (r *documentsRepository) List(ctx context.Context) ([]entity.Document, error) {
	var rows []model.Document
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]entity.Document, 0, len(rows))
	for _, m := range rows {
		result = append(result, toEntity(m))
	}
	return result, nil
}

func (r *documentsRepository) Update(ctx context.Context, doc entity.Document) (entity.Document, error) {
	updates := map[string]any{
		"request_id":       doc.RequestID,
		"template_code":    doc.TemplateCode,
		"template_version": doc.TemplateVersion,
		"payload":          doc.Payload,
		"metadata":         doc.Metadata,
		"status":           doc.Status,
		"error_message":    doc.ErrorMessage,
		"file_name":        doc.FileName,
		"file_path":        doc.FilePath,
		"file_size":        doc.FileSize,
		"checksum":         doc.Checksum,
		"content_type":     doc.ContentType,
		"store_to_dms":     doc.StoreToDMS,
		"dms_document_id":  doc.DMSDocumentID,
		"dms_status":       doc.DMSStatus,
		"has_callback":     doc.HasCallback,
		"callback_url":     doc.CallbackURL,
		"callback_status":  doc.CallbackStatus,
		"created_by":       doc.CreatedBy,
		"processed_at":     doc.ProcessedAt,
		"updated_at":       time.Now(),
	}
	tx := r.db.WithContext(ctx).Model(&model.Document{}).Where("id = ?", doc.ID).Updates(updates)
	if tx.Error != nil {
		return entity.Document{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entity.Document{}, errors.New("document not found")
	}
	return doc, nil
}

func (r *documentsRepository) Delete(ctx context.Context, id int64) error {
	tx := r.db.WithContext(ctx).Delete(&model.Document{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("document not found")
	}
	return nil
}

func toModel(e entity.Document) model.Document {
	return model.Document{
		ID:              e.ID,
		RequestID:       e.RequestID,
		TemplateCode:    e.TemplateCode,
		TemplateVersion: e.TemplateVersion,
		Payload:         e.Payload,
		Metadata:        e.Metadata,
		Status:          e.Status,
		ErrorMessage:    e.ErrorMessage,
		FileName:        e.FileName,
		FilePath:        e.FilePath,
		FileSize:        e.FileSize,
		Checksum:        e.Checksum,
		ContentType:     e.ContentType,
		StoreToDMS:      e.StoreToDMS,
		DMSDocumentID:   e.DMSDocumentID,
		DMSStatus:       e.DMSStatus,
		HasCallback:     e.HasCallback,
		CallbackURL:     e.CallbackURL,
		CallbackStatus:  e.CallbackStatus,
		CreatedBy:       e.CreatedBy,
		CreatedAt:       e.CreatedAt,
		ProcessedAt:     e.ProcessedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func toEntity(m model.Document) entity.Document {
	return entity.Document{
		ID:              m.ID,
		RequestID:       m.RequestID,
		TemplateCode:    m.TemplateCode,
		TemplateVersion: m.TemplateVersion,
		Payload:         m.Payload,
		Metadata:        m.Metadata,
		Status:          m.Status,
		ErrorMessage:    m.ErrorMessage,
		FileName:        m.FileName,
		FilePath:        m.FilePath,
		FileSize:        m.FileSize,
		Checksum:        m.Checksum,
		ContentType:     m.ContentType,
		StoreToDMS:      m.StoreToDMS,
		DMSDocumentID:   m.DMSDocumentID,
		DMSStatus:       m.DMSStatus,
		HasCallback:     m.HasCallback,
		CallbackURL:     m.CallbackURL,
		CallbackStatus:  m.CallbackStatus,
		CreatedBy:       m.CreatedBy,
		CreatedAt:       m.CreatedAt,
		ProcessedAt:     m.ProcessedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

