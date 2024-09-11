package endpoints

import (
	"net/http"
	"sg-business-service/modules/outstanding/followups"
	"sg-business-service/utils"
)

func SampleFollowUp(res http.ResponseWriter, req  *http.Request) {

	followup := &followups.FollowUpCreationRequest{
		followup: followups.FollowUp{
			ID:                [12]byte{},
			RefPrevFollowUpId: new(string),
			FollowUpId:        "",
			ContactPersonId:   "",
			PersonInChargeId:  0,
			PartyName:         "",
			Description:       "",
			Status:            0,
			FollowUpBills:     []followups.FollowUpBill{},
		},		
	}

	response := utils.NewResponseStruct(followup, 1)
	response.ToJson(res)
}