package ledgers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sg-business-service/config"
	"sg-business-service/handlers"
	"sg-business-service/models"
	"sg-business-service/utils"
)

func getCollection() *handlers.MongoHandler {
	var collection = handlers.GetCollection(config.TallyDb, config.Ledger)
	return handlers.NewMongoHandler(collection)
}

func GetLedgersByNames(companyId string, partyNames []string) []MetaLedger {
	var filter = bson.M{
		"GUID": bson.M{
			"$regex": "^" + companyId,
		},
		"Name": bson.M{
			"$in": partyNames,
		},
	}

	docFilter := handlers.DocumentFilter{
		Ctx:           context.TODO(),
		Filter:        filter,
		UsePagination: false,
		Limit:         0,
		Offset:        0,
		Projection: bson.M{
			"Name":         1,
			"Group":        1,
			"Address":      1,
			"State":        1,
			"PinCode":      1,
			"Email":        1,
			"EmailCc":      1,
			"CreditPeriod": 1,
			"CreditLimit":  1,
			"_id":          0,
		},
		Sorting: bson.D{
			{
				Key:   "Name",
				Value: utils.GetValueBySortOrder("asc"),
			},
		},
	}

	var handler = getCollection()
	ledgers, err := handlers.GetDocuments[MetaLedger](handler, docFilter)
	if err != nil {
		log.Fatal(err)
		return make([]MetaLedger, 0)
	}
	return ledgers
}

func GetLedgers(companyId string, requestFilter models.RequestFilter, additionalFilter []bson.M) []MetaLedger {
	var filter = bson.M{
		"GUID": bson.M{
			"$regex": "^" + companyId,
		},
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
			"Name":         1,
			"Group":        1,
			"Address":      1,
			"State":        1,
			"PinCode":      1,
			"Email":        1,
			"EmailCc":      1,
			"CreditPeriod": 1,
			"CreditLimit":  1,
			"_id":          0,
		},
		Sorting: bson.D{
			{
				Key:   requestFilter.SortKey,
				Value: utils.GetValueBySortOrder(requestFilter.SortOrder),
			},
		},
	}

	if len(requestFilter.SortKey) < 1 {
		docFilter.Sorting[0].Key = "Name"
	}

	var handler = getCollection()
	ledgers, err := handlers.GetDocuments[MetaLedger](handler, docFilter)
	if err != nil {
		log.Fatal(err)
		return make([]MetaLedger, 0)
	}
	return ledgers
}

func GetLedgersByPincodes(companyId string, requestFilter models.RequestFilter, pincodes *[]string) []MetaLedger {
	var collectionFilter bson.M = nil
	if pincodes != nil {
		collectionFilter = bson.M{
			"PinCode": bson.M{
				"$in": pincodes,
			},
		}
	}
	return GetLedgers(companyId, requestFilter, []bson.M{collectionFilter})

}

func GetLedgersByStates(companyId string, states []string, requestFilter models.RequestFilter) []MetaLedger {
	collectionFilter := bson.M{
		"State": bson.M{
			"$in": states,
		},
	}
	return GetLedgers(companyId, requestFilter, []bson.M{collectionFilter})
}
