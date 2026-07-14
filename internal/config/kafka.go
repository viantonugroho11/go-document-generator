package config

type Kafka struct {
	Brokers   []string `json:"brokers"`
	ClientID  string   `json:"client_id"`
	GroupID   string   `json:"group_id"`
	Topic     string   `json:"topic"`
	TopicTemplates        string `json:"topic_templates"`
	TopicTemplateVersions string `json:"topic_template_versions"`
	TopicDocuments        string `json:"topic_documents"`
	// TopicDocumentProcess adalah dedicated processing queue untuk generation consumer.
	// Consumer subscribe ke topic ini, bukan document-events (observability).
	// Default: "document-process"
	TopicDocumentProcess  string `json:"topic_document_process"`
	GroupIDDocumentWorker string `json:"group_id_document_worker"`
	// Consumer kedua (contoh: order)
	TopicOrders   string `json:"topic_orders"`
	GroupIDOrders string `json:"group_id_orders"`
}