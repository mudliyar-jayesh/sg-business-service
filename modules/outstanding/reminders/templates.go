package reminders

import (
	"context"
	"sg-business-service/config"
	"sg-business-service/handlers"

	"go.mongodb.org/mongo-driver/bson"
)

func GetTemplateCollection() *handlers.MongoHandler {
	settingCollection := handlers.GetCollection(config.AppDb, config.OsTemplate)
	return handlers.NewMongoHandler(settingCollection)
}

func Create(template OutstandingTemplate) {
	handlers.InsertDocument(config.AppDb, config.OsTemplate, template)
}

func GetByTemplateName(companyId, templateName string) *OutstandingTemplate {
	docFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		UsePagination: false,
		Filter: bson.M{
			"company_id":    companyId,
			"template_name": templateName,
		},
	}
	templates, err := handlers.GetDocuments[OutstandingTemplate](GetTemplateCollection(), docFilter)
	if err != nil || len(templates) < 1 {
		return nil
	}
	return &templates[0]
}

func Get(companyId string) []OutstandingTemplate {
	docFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		UsePagination: false,
		Filter: bson.M{
			"company_id": companyId,
		},
	}
	templates, err := handlers.GetDocuments[OutstandingTemplate](GetTemplateCollection(), docFilter)
	if err != nil {
		return make([]OutstandingTemplate, 0)
	}
	return templates
}
