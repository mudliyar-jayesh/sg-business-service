package collection

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"sg-business-service/config"
	"sg-business-service/handlers"
	"sg-business-service/models"
	"sg-business-service/utils"
)

func getFieldBySortKey(sortKey string) string {
	var fieldBySortKey = make(map[string]string)
	fieldBySortKey["Party"] = "LedgerName"
	fieldBySortKey["Group"] = "LedgerGroupName"
	fieldBySortKey["Bill"] = "Name"

	sortField, exists := fieldBySortKey[sortKey]
	if exists {
		return sortField
	}
	return fieldBySortKey["Party"]
}
func GetFieldBySearchKey(searchKey string) string {
	var fieldBySearchKey = make(map[string]string)
	fieldBySearchKey["Party"] = "LedgerName"
	fieldBySearchKey["Group"] = "LedgerGroupName"
	fieldBySearchKey["Bill"] = "Name"

	searchField, exists := fieldBySearchKey[searchKey]
	if exists {
		return searchField
	}
	return fieldBySearchKey["Party"]
}

func GetCollection() *handlers.MongoHandler {
	var collection = handlers.GetCollection(config.TallyDb, config.Bill)
	return handlers.NewMongoHandler(collection)
}

func GetOverViewByFilter(mongoFilter bson.M, filter models.RequestFilter) []CollectionOverview {
	docFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		Filter:        mongoFilter,
		UsePagination: filter.Batch.Apply,
		Limit:         filter.Batch.Limit,
		Offset:        filter.Batch.Offset,
		Projection: bson.M{
			"LedgerName":      1,
			"LedgerGroupName": 1,
			"PendingAmount":   "$ClosingBal.Amount",
			"OpeningAmount":   "$OpeningBal.Amount",
			"BillDate":        "$BillDate.Date",
			"DueDate":         "$BillCreditPeriod.DueDate",
			"Name":            1,
			"_id":             0,
		},
		Sorting: bson.D{
			{
				Key:   getFieldBySortKey(filter.SortKey),
				Value: utils.GetValueBySortOrder(filter.SortOrder),
			},
		},
	}

	var handler = GetCollection()
	collections, err := handlers.GetDocuments[CollectionOverview](handler, docFilter)

	if err != nil {
		fmt.Println("Error occured in while getting collection")
		return make([]CollectionOverview, 0)
	}
	return collections
}
