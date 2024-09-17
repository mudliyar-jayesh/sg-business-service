package endpoints

import (
	"net/http"
	"sg-business-service/models"
	osMod "sg-business-service/modules/outstanding"
	"sg-business-service/utils"
)

func GetPartySummary(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	var body, reqErr = utils.ReadRequestBody[models.RequestFilter](req)
	if reqErr != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	filter := utils.GenerateSearchFilter(body.SearchText, "LedgerName")
	requestFilter := models.RequestFilter{}

	var bills = osMod.GetBills(headers.CompanyId, requestFilter, true, filter)

	var partyBills = utils.GroupFor(bills, func(entry osMod.MetaBill) string {
		return entry.PartyName
	})

	var partyOverview []osMod.PartyOverview
	for key, bills := range partyBills {
		var overview = osMod.PartyOverview{
			PartyName:    key,
			TotalBills:   int32(len(bills)),
			TotalOpening: 0,
			TotalClosing: 0,
		}

		for _, bill := range bills {

			if bill.OpeningAmount != nil {
				overview.TotalOpening += bill.OpeningAmount.Value
			}
			if bill.PendingAmount != nil {
				overview.TotalClosing += bill.PendingAmount.Value
			}
		}
		partyOverview = append(partyOverview, overview)
	}

	var offset = body.Batch.Offset
	var limit = body.Batch.Limit
	var billLength = int64(len(partyOverview))

	if offset >= billLength {
		response := utils.NewResponseStruct(make([]osMod.PartyOverview, 0), 0)
		response.ToJson(res)
		return
	}

	end := offset + limit
	if end > billLength {
		end = billLength // Adjust the end if it exceeds the list size
	}

	partyOverview = partyOverview[offset:end]
	response := utils.NewResponseStruct(partyOverview, len(partyOverview))
	response.ToJson(res)
}
