package postgres

import (
	"context"

	cbmodel "go-document-generator/internal/repository/documentcallbackattempts/model"
	logmodel "go-document-generator/internal/repository/documentrenderlogs/model"
	tplmodel "go-document-generator/internal/repository/documenttemplates/model"
	vermodel "go-document-generator/internal/repository/documenttemplateversions/model"
	docmodel "go-document-generator/internal/repository/documents/model"
	"go-document-generator/internal/repository/user/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(ctx context.Context, dsn string) (*gorm.DB, error) {
	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	}
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&tplmodel.DocumentTemplate{},
		&vermodel.DocumentTemplateVersion{},
		&docmodel.Document{},
		&logmodel.DocumentRenderLog{},
		&cbmodel.DocumentCallbackAttempt{},
	)
}
