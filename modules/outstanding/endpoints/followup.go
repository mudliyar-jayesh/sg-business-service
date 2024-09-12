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
	listOfStatus := [...]string{"Pending", "Scheduled", "Completed"}

	response := utils.NewResponseStruct(listOfStatus, 3)
	response.ToJson(res)
}

func CreateFollowUp(res http.ResponseWriter, req *http.Request){
	headers, err := utils.ResolveHeaders(&req.Header)

	if err != nil {
		// Give error
	}

	requestBody, err := utils.ReadRequestBody[followups.FollowUpCreationRequest](req)

	requestBody.Followup.PersonInChargeId = headers.UserId

	if err != nil {
		fmt.Printf("Error while parsing request body %v", err)
	}

	followups.CreateFollowUp(requestBody.Followup, &requestBody.PointOfContact)

}