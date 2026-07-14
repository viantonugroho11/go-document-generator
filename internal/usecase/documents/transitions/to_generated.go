package transitions

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	docEntity "go-document-generator/internal/entity/documents"
	"go-document-generator/internal/entity/enums"
	"go-document-generator/internal/shared/storage"
)

type toGenerated struct {
	deps Deps
}

func NewToGenerated(deps Deps) *toGenerated {
	return &toGenerated{deps: deps}
}

func (h *toGenerated) OnStateTransition(ctx context.Context, update docEntity.Document) (docEntity.Document, error) {
	d := update
	if err := generateAndFinalize(ctx, h.deps, &d); err != nil {
		msg := err.Error()
		d.Status = enums.DocumentStatusFailed
		d.ErrorMessage = &msg
		return d, err
	}
	return d, nil
}

func generateAndFinalize(ctx context.Context, deps Deps, d *docEntity.Document) error {
	if deps.Selector == nil {
		return errors.New("document generator not configured")
	}
	if d.TemplateID == nil || d.TemplateVersionID == nil {
		return errors.New("document template reference is missing")
	}

	tpl, err := deps.Templates.GetByID(ctx, nil, *d.TemplateID, d.TenantID)
	if err != nil {
		return err
	}
	ver, err := deps.Versions.GetByID(ctx, nil, *d.TemplateID, *d.TemplateVersionID, d.TenantID)
	if err != nil {
		return err
	}

	gen := deps.Selector.Select(string(d.OutputFormat), string(tpl.Engine))
	data, contentType, err := gen.Generate(ctx, ver.Content, d.Payload)
	if err != nil {
		return fmt.Errorf("generate document: %w", err)
	}

	ext := storage.ExtensionForFormat(string(d.OutputFormat))

	var path, fileName string
	var storageProvider enums.StorageProvider
	if deps.Storage != nil {
		path, fileName, err = deps.Storage.Save(ctx, d.ID, d.RequestID, ext, data)
		if err != nil {
			return fmt.Errorf("save document file: %w", err)
		}
		storageProvider = enums.StorageProviderMinio
	} else {
		path, fileName, err = storage.SaveDocument("", d.ID, d.RequestID, ext, data)
		if err != nil {
			return fmt.Errorf("save document file: %w", err)
		}
		storageProvider = enums.StorageProviderLocal
	}

	sum := sha256.Sum256(data)
	chk := hex.EncodeToString(sum[:])
	size := int64(len(data))
	now := time.Now().UTC()
	provider := storageProvider

	d.Status = enums.DocumentStatusGenerated
	d.FilePath = &path
	d.FileName = &fileName
	d.ContentType = &contentType
	d.FileSize = &size
	d.Checksum = &chk
	d.StorageProvider = &provider
	d.ProcessedAt = &now
	d.ErrorMessage = nil

	return nil
}
