package overview

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"reflect"
	"sg-business-service/config"
	"sg-business-service/handlers"
	"sg-business-service/utils"
	"sort"
	"time"
)

func getParentGroups(companyId string, isDebit bool) []string {
	var groupType = "Current Liabilities"
	if isDebit {
		groupType = "Current Assets"
	}

	return handlers.CachedGroups.GetChildrenNames(companyId, groupType)
}

func getCollection() *handlers.MongoHandler {
	var collection = handlers.GetCollection(config.TallyDb, config.Bill)
	return handlers.NewMongoHandler(collection)
}

func getBills(companyId string, filter OverviewFilter, additionalFilter *[]bson.M) []Bill {
	var dbFilter = bson.M{
		"CompanyId": companyId,
		"LedgerGroupName": bson.M{
			"$in": filter.Groups,
		},
		"ClosingBal":         bson.M{"$ne": nil},
		"ClosingBal.IsDebit": filter.IsDebit,
	}
	if additionalFilter != nil {
		dbFilter["$and"] = additionalFilter
	}

	docFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		Filter:        dbFilter,
		UsePagination: filter.Filter.Batch.Apply,
		Limit:         filter.Filter.Batch.Limit,
		Offset:        filter.Filter.Batch.Offset,
		Projection: bson.M{
			"Name":            1,
			"LedgerName":      1,
			"LedgerGroupName": 1,
			"BillDate":        "$BillDate.Date",
			"DueDate":         "$BillCreditPeriod.DueDate",
			"ClosingBalance":  "$ClosingBal.Amount",
			"OpeningBalance":  "$OpeningBal.Amount",
			"IsAdvance":       "$IsAdvance.Value",
			"_id":             0,
		},
		Sorting: bson.D{
			{
				Key:   filter.Filter.SortKey,
				Value: utils.GetValueBySortOrder(filter.Filter.SortOrder),
			},
		},
	}
	var handler = getCollection()
	bills, err := handlers.GetDocuments[Bill](handler, docFilter)
	if err != nil {
		log.Fatal(err)
		return make([]Bill, 0)
	}
	return bills

}
