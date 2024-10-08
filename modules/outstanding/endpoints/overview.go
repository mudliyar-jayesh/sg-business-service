package endpoints

import (
	"net/http"
	"sg-business-service/modules/outstanding/overview"
	"sg-business-service/utils"
)

func GetAgingOverview(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")
	applyRange := utils.GetBoolFromQuery(req, "applyRange")

	reqBody, err := utils.ReadRequestBody[overview.OverviewFilter](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}

  var partyOverview = overview.GetAgingOverview(companyId, applyRange, *reqBody)
	response := utils.NewResponseStruct(partyOverview, len(partyOverview))
	response.ToJson(res)

}
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

func GetBillOverview(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")

	reqBody, err := utils.ReadRequestBody[overview.OverviewFilter](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var billOverview = overview.GetBillWiseOverview(companyId, *reqBody)

	response := utils.NewResponseStruct(billOverview, len(billOverview))
	response.ToJson(res)
}
