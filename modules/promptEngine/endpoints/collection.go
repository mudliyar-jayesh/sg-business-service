package endpoints

import (
	"net/http"
	promptEngine "sg-business-service/modules/promptEngine"
	"sg-business-service/modules/promptEngine/prompts"
	"sg-business-service/utils"
)

func GetCollectionPrompts(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	body, err := utils.ReadRequestBody[promptEngine.PromptRequest](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
	}

	var userPrompts = prompts.GetCollectionPrompt(headers.CompanyId, body.StartDateStr, body.EndDateStr, body.Filter)

	response := utils.NewResponseStruct(userPrompts, len(userPrompts))
	response.ToJson(res)
}
