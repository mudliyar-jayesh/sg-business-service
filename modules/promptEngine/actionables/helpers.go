package actionables

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"sg-business-service/config"
	"sg-business-service/handlers"
	promptMod "sg-business-service/modules/promptEngine"
)

func GetCollection() *handlers.MongoHandler {
	var collection = handlers.GetCollection(config.SummaryDb, config.CollectionActionables)
	return handlers.NewMongoHandler(collection)
}

func InsertOne(companyId string, actionable promptMod.Actionable) string {
	guid := uuid.New()
	actionable.Guid = guid.String()
	actionable.CompanyId = companyId
	_, err := handlers.InsertDocument[promptMod.Actionable](config.SummaryDb, config.CollectionActionables, actionable)
	if err != nil {
		return ""
	}
	return actionable.Guid
}

func GetByGuid(companyId, guid string) *promptMod.Actionable {
	handler := GetCollection()
	filter := bson.M{
		"guid":       guid,
		"company_id": companyId,
	}
	docFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		UsePagination: true,
		Limit:         1,
		Offset:        0,
		Filter:        filter,
	}
	values, err := handlers.GetDocuments[promptMod.Actionable](handler, docFilter)
	if err != nil || len(values) < 1 {
		return nil
	}
	return &values[0]

}

func UpdateOne(actionable promptMod.Actionable) {
	handler := GetCollection()
	filter := bson.M{
		"guid": actionable.Guid,
	}
	handler.ReplaceDocument(config.SummaryDb, config.CollectionActionables, filter, actionable)
}

func GetByCompanyId(companyId string, additionalFilter *bson.M) []promptMod.Actionable {
	handler := GetCollection()
	filter := bson.M{
		"company_id": companyId,
	}
	docFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		UsePagination: false,
		Filter:        filter,
	}
	values, err := handlers.GetDocuments[promptMod.Actionable](handler, docFilter)
	if err != nil || len(values) < 1 {
		return make([]promptMod.Actionable, 0)
	}
	return values

}
