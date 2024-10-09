package endpoints

import (
	"fmt"
	"net/http"
	"sg-business-service/config"
	"sg-business-service/modules/outstanding/reminders"
	promptEngine "sg-business-service/modules/promptEngine"
	"sg-business-service/modules/promptEngine/actionables"
	"sg-business-service/modules/promptEngine/prompts"
	"sg-business-service/utils"
	"time"
)

func GetCollectionPrompts(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	var partyWise = utils.GetBoolFromQuery(req, "partyWise")

	body, err := utils.ReadRequestBody[promptEngine.PromptRequest](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
	}

	var userPrompts = prompts.GetCollectionPrompt(headers.CompanyId, body.StartDateStr, body.EndDateStr, body.Filter, partyWise)

	response := utils.NewResponseStruct(userPrompts, len(userPrompts))
	response.ToJson(res)
}

func ProcessCollectionDecision(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	var partyWise = utils.GetBoolFromQuery(req, "partyWise")

	body, err := utils.ReadRequestBody[promptEngine.DecisionRequest](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
	}

	switch body.ActionCode {
	case prompts.SendReminderAction:
		var parties = make([]string, 1)
		parties[0] = *body.PartyName
		reminders.SendEmailReminder(headers.CompanyId, parties)
	case prompts.TeamFollowUpAction:
		message := fmt.Sprintf("Follow up with %v for total amount of %v", *body.PartyName, *body.AmountStr)
		if !partyWise {
			message = fmt.Sprintf("Follow up with %v for bill number: %v with amount of %v", *body.PartyName, *body.BillNumber, *body.AmountStr)
		}
		var actionable = promptEngine.Actionable{
			Title:       "Collection Follow-Up Task",
			Description: message,
			AssignedTo:  0,
			CreatedBy:   int64(headers.UserId),
			CreatedOn:   time.Now().UTC(),
			Status:      promptEngine.Pending,
		}
		actionables.InsertOne(headers.CompanyId, actionable)
	}
}

func GetCollectionActionables(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	var url string = config.PortalUrl + "/companies/get/users"
	users := utils.GetFromPortal[[]config.MetaUser](url, headers)
	userById := utils.ToDict(*users, func(user config.MetaUser) int64 {
		return int64(user.ID)
	})

	var tasks = actionables.GetByCompanyId(headers.CompanyId, nil)

	var collectionTasks []promptEngine.ActionableOverview
	for _, task := range tasks {
		user, exists := userById[task.AssignedTo]
		assignedTo := "Other"
		if exists {
			assignedTo = user.Name
		}
		if task.AssignedTo == int64(headers.UserId) {
			assignedTo = "Self"
		}

		user, exists = userById[task.CreatedBy]
		createdBy := "Other"
		if exists {
			createdBy = user.Name
		}
		if task.CreatedBy == int64(headers.UserId) {
			createdBy = "Self"
		}
		var status string = "Pending"
		if task.Status == promptEngine.Cancelled {
			status = "Cancelled"
		} else if task.Status == promptEngine.Done {
			status = "Done"
		}

		var collectionTask = promptEngine.ActionableOverview{
			Task:           task,
			AssignedToName: assignedTo,
			CreatedByName:  createdBy,
			StatusName:     status,
		}
		collectionTasks = append(collectionTasks, collectionTask)
	}

	response := utils.NewResponseStruct(collectionTasks, len(collectionTasks))
	response.ToJson(res)
}

func UpdateCollectionTask(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}
	body, err := utils.ReadRequestBody[promptEngine.Actionable](req)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
	}

	actionables.UpdateOne(*body)
}
