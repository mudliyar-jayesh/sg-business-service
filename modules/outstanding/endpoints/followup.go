package endpoints

import (
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