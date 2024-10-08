package endpoints

import (
  "net/http"
  "sg-business-service/utils"
  vchMod "sg-business-service/modules/vouchers"
)

func GetVoucherInfo(w http.ResponseWriter, r *http.Request) {
	headers, err := utils.ResolveHeaders(&r.Header)
	if headers.HandleErrorOrIllegalValues(w, &err) {
		return
	}
  var voucherIds []string
  var ledgerNames []string

  var voucherTypes []string
  voucherTypes = append(voucherTypes, "K - Sales GST")

  var vouchers = vchMod.GetMetaVouchers(headers.CompanyId, voucherTypes, voucherIds, ledgerNames)
  response := utils.NewResponseStruct(vouchers, len(vouchers))
  response.ToJson(w)
}

