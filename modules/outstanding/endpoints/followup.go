package endpoints

import (
	"fmt"
	"net/http"
	"sg-business-service/models"
	osMod "sg-business-service/modules/outstanding"
	"sg-business-service/modules/outstanding/followups"
	"sg-business-service/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func SampleFollowUp(res http.ResponseWriter, req  *http.Request) {
	followup := &followups.FollowUpCreationRequest{
	}
	response := utils.NewResponseStruct(followup, 1)
	response.ToJson(res)
}

func GetBillStatusList(res http.ResponseWriter, req *http.Request){
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err){return}
		
	mappings := followups.GetFollowUpStatusMappings()
	response := utils.NewResponseStruct(mappings, 1)
	response.ToJson(res)
}

func GetBills(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err){return}

	partyName := req.URL.Query().Get("partyName")
	searchText := req.URL.Query().Get("searchText")

	reqFilter := models.RequestFilter{Batch: models.Pagination{Apply: true, Limit: 25}}

	var filter = []bson.M {
		{
			"LedgerName": partyName,
		},
	}

	if len(searchText) > 0 {
		filter = append(filter, utils.GenerateSearchFilter(searchText, "Name")[0])
	}
	
	bills := osMod.GetBills(headers.CompanyId, reqFilter, true, filter)

	response := utils.NewResponseStruct(bills, len(bills))
	response.ToJson(res)
}

func GetContactPerson(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err){return}

	partyName := req.URL.Query().Get("partyName")
	contactPersons := followups.GetContactPersons(headers.CompanyId, partyName)

	response := utils.NewResponseStruct(contactPersons, len(contactPersons))
	response.ToJson(res)
}

func UpdateFollowUp(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {return}

	requestBody, err := utils.ReadRequestBody[followups.FollowUpCreationRequest](req)

	if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf("Error while parsing request body %v", err)))
			return
		}

	// TODO: change
	followups.CreateFollowUp(requestBody.Followup, nil)
}

func CreateFollowUp(res http.ResponseWriter, req *http.Request){
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err){return}
		
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Attempt to unauthorized access without secure headers"))
		return
	}

	requestBody, err := utils.ReadRequestBody[followups.FollowUpCreationRequest](req)

	if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf("Error while parsing request body %v", err)))
			return
		}

	requestBody.Followup.PersonInChargeId = headers.UserId
	requestBody.Followup.CompanyId = headers.CompanyId

	err = followups.CreateFollowUp(requestBody.Followup, &requestBody.PointOfContact)

	if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf("Error while creating followup %v", err)))
			return
		}
}

func GetFollowupList(res http.ResponseWriter,  req *http.Request){
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err){return}

	partyName := req.URL.Query().Get("partyName")

	partyFollowups := followups.GetFollowUpList(headers.CompanyId, partyName)

	response := utils.NewResponseStruct(partyFollowups, len(partyFollowups))
	response.ToJson(res)
}