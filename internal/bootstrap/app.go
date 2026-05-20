package bootstrap

import "go-document-generator/internal/transport/apis"

// RunApp memuat config (global), wiring terisolasi (DB, Redis, Kafka, routes), lalu jalankan HTTP server sampai signal.
func RunApp() error {
	if err := LoadConfig(); err != nil {
		return err
	}

	db, err := initDB()
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var cleanups []func()
	services := apis.Services{}

	userService, closeUser, err := wireUserService(db)
	if err != nil {
		return err
	}
	cleanups = append(cleanups, closeUser)
	services.Users = userService

	docServices, closeDoc, err := wireDocumentServices(db)
	if err != nil {
		for _, fn := range cleanups {
			fn()
		}
		return err
	}
	cleanups = append(cleanups, closeDoc)
	services.Templates = docServices.Templates
	services.TemplateVersions = docServices.TemplateVersions
	services.Documents = docServices.Documents
	services.RenderLogs = docServices.RenderLogs
	services.Callbacks = docServices.Callbacks

	defer func() {
		for _, fn := range cleanups {
			fn()
		}
	}()

	redisClient, err := initRedis()
	if err != nil {
		return err
	}
	defer redisClient.Close()

	e := newEcho(services)
	return runHTTP(e)
}
