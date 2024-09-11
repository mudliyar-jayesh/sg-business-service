package endpoints

import (
	"net/http"
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
