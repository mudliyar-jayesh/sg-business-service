package endpoints

import (
	"net/http"
	"sg-business-service/modules/outstanding/reminders"
	"sg-business-service/utils"
)

func CreateOsTemplates(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")

	newTemplate, err := utils.ReadRequestBody[reminders.OutstandingTemplate](req)
	newTemplate.CompanyId = companyId
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}

	reminders.Create(*newTemplate)

}

func GetAllOsTemplates(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")
	templates := reminders.Get(companyId)
	response := utils.NewResponseStruct(templates, len(templates))
	response.ToJson(res)
}

func GetOsTemplatesByName(res http.ResponseWriter, req *http.Request) {
	companyId := req.Header.Get("CompanyId")
	name := req.URL.Query().Get("name")
	template := reminders.GetByTemplateName(companyId, name)
	response := utils.NewResponseStruct(template, 1)
	response.ToJson(res)
}
