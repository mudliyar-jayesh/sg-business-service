package endpoints

import (
    "net/http"
    "sg-business-service/modules/syncInfo"
    "sg-business-service/utils"
)

func GetLastSync(res http.ResponseWriter, req *http.Request) {
    companyId := req.Header.Get("CompanyId")

    info := syncInfo.GetByCompanyId(companyId)
    response := utils.NewResponseStruct(info, 1)
    response.ToJson(res);

}
