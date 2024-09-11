package collection

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sg-business-service/handlers"
	"sg-business-service/utils"
	"time"
)

const parentName string = "Current Assets"

func GetCollectionOverview(companyId string, collectionFilter CollectionFilter) []CollectionOverview {
	var groups = handlers.CachedGroups.GetChildrenNames(companyId, parentName)

	if len(collectionFilter.Groups) > 0 {
		groups = utils.Intersection(groups, collectionFilter.Groups)
	}
	startDate, err := time.Parse("02-01-2006", collectionFilter.StartDateStr)
	if err != nil {
		return make([]CollectionOverview, 0)
	}
	endDate, err := time.Parse("02-01-2006", collectionFilter.EndDateStr)
	if err != nil {
		return make([]CollectionOverview, 0)
	}

	var filter bson.M = bson.M{
		"CompanyId": companyId,
		"LedgerGroupName": bson.M{
			"$in": groups,
		},
		"$and": []bson.M{
			{"BillDate.Date": bson.M{"$exists": true, "$ne": nil}},
			{"BillDate.Date": bson.M{
				"$gte": primitive.NewDateTimeFromTime(startDate),
				"$lte": primitive.NewDateTimeFromTime(endDate),
			}},
		},
	}

	if len(collectionFilter.Parties) > 0 {
		filter["LedgerName"] = bson.M{
			"$in": collectionFilter.Parties,
		}
	}
	if len(collectionFilter.Filter.SearchText) > 0 {
		var searchField = GetFieldBySearchKey(collectionFilter.Filter.SearchKey)
		filter["$and"] = utils.GenerateSearchFilter(collectionFilter.Filter.SearchText, searchField)
	}
	collectionFilter.Filter.Batch.Apply = collectionFilter.Filter.Batch.Limit != 0

	return GetOverViewByFilter(filter, collectionFilter.Filter)
}
