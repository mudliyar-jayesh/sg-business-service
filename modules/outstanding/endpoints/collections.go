package endpoints

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"sg-business-service/handlers"
	"sg-business-service/modules/outstanding/collection"
	"sg-business-service/utils"
)

func GetCollectionOverview(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")

	collectionFilter, err := utils.ReadRequestBody[collection.CollectionFilter](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}

	collections := collection.GetCollectionOverview(companyId, *collectionFilter)

	response := utils.NewResponseStruct(collections, len(collections))
	response.ToJson(res)
}

func TestCollectionFetch(res http.ResponseWriter, req *http.Request) {
	docFilter := handlers.DocumentFilter{
		Ctx: context.TODO(),
		Filter: bson.M{
			"ClosingBal": bson.M{
				"$ne": nil,
			},
		},
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
	}

	values, err := handlers.GetDocuments[collection.CollectionOverview](collection.GetCollection(), docFilter)
	if err != nil {
		http.Error(res, "Unable to read data ", http.StatusBadRequest)
		return
	}
	response := utils.NewResponseStruct(values, len(values))
	response.ToJson(res)
}
