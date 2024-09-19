package endpoints

import (
	"net/http"
	"sg-business-service/modules/outstanding/summary"
	"sg-business-service/utils"
)

func CalculuateSummary(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}
	summary.CalculateOutstandingSummary(headers.CompanyId)
}
