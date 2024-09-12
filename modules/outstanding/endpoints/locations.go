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
)

func GetLocationWiseOverview(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")
	isDebit := utils.GetBoolFromQuery(req, "isDebit")

	body, err := utils.ReadRequestBody[osMod.OsLocationFilter](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}

	requestFilter := models.RequestFilter{
		Batch: models.Pagination{
			Apply: false,
		},
	}

	var ledgers []ledgerMod.MetaLedger

	switch body.LocationType {
	case osMod.StateWise:
		stateNames := dsp.GetStates()
		ledgers = ledgerMod.GetLedgersByStates(companyId, stateNames, requestFilter)
	case osMod.PincodeWise:
		ledgers = ledgerMod.GetLedgersByPincodes(companyId, requestFilter, nil)
	case osMod.RegionWise:
		regions := dsp.GetRegions(body.State)
		dspEntries := dsp.GetByRegions(regions)
		var pincodes []string
		for _, value := range dspEntries {
			pincodes = append(pincodes, value.Pincode)
		}
		ledgers = ledgerMod.GetLedgersByPincodes(companyId, requestFilter, &pincodes)
	case osMod.DistrictWise:
		districts := dsp.GetDistrictsByState(body.State)
		dspEntries := dsp.GetByDistrict(districts)
		var pincodes []string
		for _, value := range dspEntries {
			pincodes = append(pincodes, value.Pincode)
		}
		ledgers = ledgerMod.GetLedgersByPincodes(companyId, requestFilter, &pincodes)
	}

	var ledgerNames []string
	for _, value := range ledgers {
		ledgerNames = append(ledgerNames, value.Name)
	}

	var allBills []osMod.MetaBill
	utils.ProcessBatch(ledgerNames, 1000, func(names []string) {

		filter := []bson.M{
			{
				"LedgerName": bson.M{
					"$in": names,
				},
			},
		}

		bills := osMod.GetBills(companyId, body.Filter, isDebit, filter)
		allBills = append(allBills, bills...)
	})

	response := utils.NewResponseStruct(allBills, len(allBills))
	response.ToJson(res)
}
