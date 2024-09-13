package endpoints

import (
	//"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"sg-business-service/models"
	"sg-business-service/modules/dsp"
	ledgerMod "sg-business-service/modules/ledgers"
	osMod "sg-business-service/modules/outstanding"
	"sg-business-service/utils"
	"sync"
)

func GetLocationWiseOverview(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")
	isDebit := utils.GetBoolFromQuery(req, "isDebit")

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
	var dspEntries = dsp.GetByStates(states)
	var dspByPincode = utils.ToLookup(dspEntries, func(entry dsp.DSP) string {
		return entry.Pincode
	})
	var pincodesUnderState []string = utils.Select(dspEntries, func(location dsp.DSP) string {
		return location.Pincode
	})

	var ledgerFilter = []bson.M{
		{
			"PinCode": bson.M{
				"$in": utils.Distinct(pincodesUnderState),
			},
		},
	}
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

	var waitGroup sync.WaitGroup
	var mutex sync.Mutex

	var locationBills = make([]osMod.LocationOverview, 0)
	for pinCode, ledgers := range ledgerNamesByPinCode {
		waitGroup.Add(1)

		go func(pinCode string, ledgers []string) {
			defer waitGroup.Done()
			var ledgerFilter = []bson.M{
				{
					"LedgerName": bson.M{
						"$in": ledgers,
					},
				},
			}
			var bills = osMod.GetBills(companyId, requestFilter, isDebit, ledgerFilter)

			var totalOpeningAmount float64 = 0
			var totalClosingAmount float64 = 0
			for _, bill := range bills {
				if bill.OpeningAmount != nil {
					totalOpeningAmount += bill.OpeningAmount.Value

					var closing = bill.OpeningAmount.Value
					if bill.PendingAmount != nil {
						closing = bill.PendingAmount.Value
					}
					totalClosingAmount += closing
				}
			}

			var locationName = pinCode

			dspValue, exists := dspByPincode[pinCode]
			if !exists {
				locationName = "Other"
			}

			if len(dspValue) > 0 {
				if body.LocationType == osMod.RegionWise {
					locationName = dspValue[0].RegionName
				} else if body.LocationType == osMod.DistrictWise {
					locationName = dspValue[0].District
				}
			}
			locationBill := osMod.LocationOverview{
				OpeningAmount: totalOpeningAmount,
				ClosingAmount: totalClosingAmount,
				LocationName:  locationName,
			}

			mutex.Lock()
			locationBills = append(locationBills, locationBill)
			mutex.Unlock()
		}(pinCode, ledgers)
	}
	waitGroup.Wait()
	groupedBills := utils.GroupByKey(locationBills, "LocationName")

	var locationOverview []osMod.LocationOverview
	for location, bills := range groupedBills {

		var totalOpening float64 = 0
		var totalClosing float64 = 0
		for _, bill := range bills {
			totalOpening += bill.OpeningAmount
			totalClosing += bill.ClosingAmount
		}

		entry := osMod.LocationOverview{
			LocationName:  location,
			OpeningAmount: totalOpening,
			ClosingAmount: totalClosing,
		}
		locationOverview = append(locationOverview, entry)
	}

	// Step 6: Create the response and send it as JSON
	response := utils.NewResponseStruct(locationOverview, len(locationOverview))
	response.ToJson(res)
}
