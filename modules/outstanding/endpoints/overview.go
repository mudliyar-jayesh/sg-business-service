package endpoints

import (
	"net/http"
	"sg-business-service/modules/outstanding/overview"
	"sg-business-service/utils"
)

func GetPartyOverview(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")

	reqBody, err := utils.ReadRequestBody[overview.OverviewFilter](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var partyOverview = overview.GetPartyWiseOverview(companyId, *reqBody)

	response := utils.NewResponseStruct(partyOverview, len(partyOverview))
	response.ToJson(res)
}
