package outstanding

import (
	"context"
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

func GetOutstandingCollection() *handlers.MongoHandler {
	var collection = handlers.GetCollection(config.TallyDb, config.Bill)
	return handlers.NewMongoHandler(collection)
}

func GetOutstandingByFilter(mongoFilter bson.M, requestFilter OsReportFilter, usePagination bool) handlers.DocumentResponse {
	docFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		Filter:        mongoFilter,
		UsePagination: usePagination,
		Limit:         requestFilter.Limit,
		Offset:        requestFilter.Offset,
		Projection: bson.M{
			"LedgerName":      1,
			"LedgerGroupName": 1,
			"BillDate":        "$BillDate.Date",
			"DueDate":         "$BillCreditPeriod.DueDate",
			"Amount":          "$ClosingBal.Amount",
			"OpeningAmount":   "$OpeningBal.Amount",
			"Name":            "$Name",
			"_id":             0,
		},
		Sorting: bson.D{
			{
				Key:   getFieldBySortKey(requestFilter.SortKey),
				Value: utils.GetValueBySortOrder(requestFilter.SortOrder),
			},
		},
	}

	var handler = GetOutstandingCollection()
	return handler.FindDocuments(docFilter)
}

func GetBills(companyId string, requestFilter models.RequestFilter, isDebit bool, additionalFilter []bson.M) []MetaBill {
	var parentName string = "Current Assets"
	if !isDebit {
		parentName = "Current Liabilities"
	}
	var groups = handlers.CachedGroups.GetChildrenNames(companyId, parentName)

	var filter = bson.M{
		"CompanyId": companyId,
		"LedgerGroupName": bson.M{
			"$in": groups,
		},
		"ClosingBal":         bson.M{"$ne": nil},
		"ClosingBal.IsDebit": isDebit,
	}
	if additionalFilter != nil {
		filter["$and"] = additionalFilter
	}

	docFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		Filter:        filter,
		UsePagination: requestFilter.Batch.Apply,
		Limit:         requestFilter.Batch.Limit,
		Offset:        requestFilter.Batch.Offset,
		Projection: bson.M{
			"LedgerName":      1,
			"LedgerGroupName": 1,
			"BillDate":        "$BillDate.Date",
			"DueDate":         "$BillCreditPeriod.DueDate",
			"Amount":          "$ClosingBal.Amount",
			"OpeningAmount":   "$OpeningBal.Amount",
			"Name":            "$Name",
			"_id":             0,
		},
		Sorting: bson.D{
			{
				Key:   getFieldBySortKey(requestFilter.SortKey),
				Value: utils.GetValueBySortOrder(requestFilter.SortOrder),
			},
		},
	}

	var handler = GetOutstandingCollection()
	bills, err := handlers.GetDocuments[MetaBill](handler, docFilter)
	if err != nil {
		return make([]MetaBill, 0)
	}
	return bills

}
