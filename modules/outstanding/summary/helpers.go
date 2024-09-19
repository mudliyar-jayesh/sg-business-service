package summary

import (
	"sg-business-service/config"
	"sg-business-service/handlers"
)

func GetCollection() *handlers.MongoHandler {
	var collection = handlers.GetCollection(config.SummaryDb, config.OsSummary)
	return handlers.NewMongoHandler(collection)
}

func ConvertToInterfaceSlice(summaries []OutstandingSummary) []interface{} {
	result := make([]interface{}, len(summaries))
	for i, v := range summaries {
		result[i] = v
	}
	return result
}
