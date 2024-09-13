package endpoints

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"sg-business-service/models"
	"sg-business-service/modules/dsp"
	ledgerMod "sg-business-service/modules/ledgers"
	osMod "sg-business-service/modules/outstanding"
	"sg-business-service/utils"
)

func GetLocationWiseOverview(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")
	//	isDebit := utils.GetBoolFromQuery(req, "isDebit")

	body, err := utils.ReadRequestBody[osMod.OsLocationFilter](req)
	if err != nil {
		http.Error(res, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Create a request filter without pagination
	requestFilter := models.RequestFilter{
		Batch: models.Pagination{
			Apply: false,
		},
	}

	states := make([]string, 1)
	states[0] = body.State
	var pincodesUnderState []string = utils.Select(dsp.GetByStates(states), func(location dsp.DSP) string {
		return location.Pincode
	})

	var ledgerFilter = []bson.M{
		{
			"PinCode": bson.M{
				"$in": utils.Distinct(pincodesUnderState),
			},
		},
	}
	fmt.Println("Pincodes: ", len(pincodesUnderState))
	// Step 1: Fetch all ledgers in a single call
	var ledgers []ledgerMod.MetaLedger = ledgerMod.GetLedgers(companyId, requestFilter, ledgerFilter)

	// Step 2: Initialize maps and slices for location details and bills processing
	ledgerNamesByPinCode := make(map[string][]string) // Map of PinCode to Ledger names

	// Step 3: Group ledgers by PinCode in a single pass
	for _, ledger := range ledgers {
		pinCode := "Other"
		if ledger.PinCode != nil {
			pinCode = *ledger.PinCode
		}
		ledgerNamesByPinCode[pinCode] = append(ledgerNamesByPinCode[pinCode], ledger.Name)
	}

	// Step 6: Create the response and send it as JSON
	response := utils.NewResponseStruct(ledgerNamesByPinCode, len(ledgerNamesByPinCode))
	response.ToJson(res)
}
