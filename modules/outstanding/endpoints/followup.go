package endpoints

import (
	"fmt"
	"net/http"
	"sg-business-service/modules/outstanding/followups"
	"sg-business-service/utils"
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



}