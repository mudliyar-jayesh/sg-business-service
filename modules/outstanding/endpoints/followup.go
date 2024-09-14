package endpoints

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sg-business-service/config"
	"sg-business-service/models"
	osMod "sg-business-service/modules/outstanding"
	fuMod "sg-business-service/modules/outstanding/followups"
	"sg-business-service/utils"
	"strconv"
	"time"
)

func SampleFollowUp(res http.ResponseWriter, req *http.Request) {
	followup := &fuMod.FollowUpCreationRequest{}
	response := utils.NewResponseStruct(followup, 1)
	response.ToJson(res)
}

func GetBillStatusList(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	mappings := fuMod.GetFollowUpStatusMappings()
	response := utils.NewResponseStruct(mappings, 1)
	response.ToJson(res)
}

func GetBills(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	partyName := req.URL.Query().Get("partyName")
	searchText := req.URL.Query().Get("searchText")

	reqFilter := models.RequestFilter{Batch: models.Pagination{Apply: true, Limit: 25}}

	var filter = []bson.M{
		{
			"LedgerName": partyName,
		},
	}

	if len(searchText) > 0 {
		filter = append(filter, utils.GenerateSearchFilter(searchText, "Name")[0])
	}

	bills := osMod.GetBills(headers.CompanyId, reqFilter, true, filter)

	response := utils.NewResponseStruct(bills, len(bills))
	response.ToJson(res)
}

func GetContactPerson(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	partyName := req.URL.Query().Get("partyName")
	contactPersons := fuMod.GetContactPersons(headers.CompanyId, partyName)

	response := utils.NewResponseStruct(contactPersons, len(contactPersons))
	response.ToJson(res)
}

func UpdateFollowUp(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	requestBody, err := utils.ReadRequestBody[fuMod.FollowUpCreationRequest](req)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Error while parsing request body %v", err)))
		return
	}

	err = fuMod.UpdateFollowUp(requestBody.Followup)

	if err != nil {
		res.WriteHeader(http.StatusExpectationFailed)
		res.Write([]byte(fmt.Sprintf("Error while updatin %v", err)))
		return
	}
}

func CreateFollowUp(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Attempt to unauthorized access without secure headers"))
		return
	}

	requestBody, err := utils.ReadRequestBody[fuMod.FollowUpCreationRequest](req)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Error while parsing request body %v", err)))
		return
	}

	requestBody.Followup.PersonInChargeId = headers.UserId
	requestBody.Followup.CompanyId = headers.CompanyId

	err = fuMod.CreateFollowUp(requestBody.Followup, requestBody.PointOfContact)
	if len(requestBody.Followup.ContactPersonId) == 0 && requestBody.PointOfContact != nil {
		requestBody.PointOfContact.CompanyId = headers.CompanyId
	}

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(fmt.Sprintf("Error while creating followup %v", err)))
		return
	}
}

func GetFollowupList(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	partyName := req.URL.Query().Get("partyName")

	partyFollowups := fuMod.GetFollowUpList(headers.CompanyId, partyName)

	response := utils.NewResponseStruct(partyFollowups, len(partyFollowups))
	response.ToJson(res)
}

func GetTeamFollowReport(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	var followups = fuMod.GetFollowups(headers.CompanyId, nil, nil)

	var followUpByMember = utils.ToLookup(followups, func(entry fuMod.FollowUp) uint64 {
		return entry.PersonInChargeId
	})

	var memberOverview []fuMod.FollowUpOverview

	var url string = config.UmsUrl + "/users/company/" + headers.CompanyId
	users := utils.GetFromUms[[]config.MetaUser](url, headers)
	userById := utils.ToDict(*users, func(user config.MetaUser) uint64 {
		return user.Id
	})

	var infoUrl string = fmt.Sprintf("%v/users/get?userId=%v", config.UmsUrl, headers.UserId)
	userInfo := utils.GetFromUms[config.MetaUser](infoUrl, headers)

	for userId, values := range followUpByMember {
		user, exists := userById[userId]
		userName := "Other"
		if exists {
			userName = user.Name
		}
		if userId == headers.UserId && userInfo != nil {
			userName = userInfo.Name
		}
		overview := fuMod.FollowUpOverview{
			Name:           userName,
			TotalCount:     0,
			PendingCount:   0,
			ScheduledCount: 0,
			CompleteCount:  0,
		}
		for _, followup := range values {
			var totalPending int32
			var totalScheduled int32
			var totalCompleted int32

			for _, bill := range followup.FollowUpBills {
				switch bill.Status {
				case fuMod.Completed:
					totalCompleted += 1
				case fuMod.Scheduled:
					totalScheduled += 1
				default:
					totalPending += 1
				}
			}

			if totalPending >= totalScheduled && totalPending >= totalCompleted {
				overview.PendingCount += 1
			} else if totalScheduled >= totalPending && totalScheduled >= totalCompleted {
				overview.ScheduledCount += 1
			} else {
				overview.CompleteCount += 1
			}
		}
		memberOverview = append(memberOverview, overview)
	}

	response := utils.NewResponseStruct(memberOverview, len(memberOverview))
	response.ToJson(res)
}

func GetPartyFollowUpReport(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	var followups = fuMod.GetFollowups(headers.CompanyId, nil, nil)

	var followUpByPartyName = utils.ToLookup(followups, func(entry fuMod.FollowUp) string {
		return entry.PartyName
	})

	var partyOverview []fuMod.FollowUpOverview

	for partyName, values := range followUpByPartyName {
		overview := fuMod.FollowUpOverview{
			Name:           partyName,
			TotalCount:     0,
			PendingCount:   0,
			ScheduledCount: 0,
			CompleteCount:  0,
		}
		for _, followup := range values {
			var totalPending int32
			var totalScheduled int32
			var totalCompleted int32

			for _, bill := range followup.FollowUpBills {
				switch bill.Status {
				case fuMod.Completed:
					totalCompleted += 1
				case fuMod.Scheduled:
					totalScheduled += 1
				default:
					totalPending += 1
				}
			}

			if totalPending >= totalScheduled && totalPending >= totalCompleted {
				overview.PendingCount += 1
			} else if totalScheduled >= totalPending && totalScheduled >= totalCompleted {
				overview.ScheduledCount += 1
			} else {
				overview.CompleteCount += 1
			}
		}
		partyOverview = append(partyOverview, overview)
	}

	response := utils.NewResponseStruct(partyOverview, len(partyOverview))
	response.ToJson(res)
}

func GetUpcomingFollowUpReport(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	body, err := utils.ReadRequestBody[fuMod.FollowUpFilter](req)
	if err != nil {
		http.Error(res, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("02-01-2006", body.StartDateStr)
	if err != nil {
		return
	}
	endDate, err := time.Parse("02-01-2006", body.EndDateStr)
	if err != nil {
		return
	}

	filter := []bson.M{
		{"CreateDate": bson.M{"$exists": true, "$ne": nil}},
		{"CreateDate": bson.M{
			"$gte": primitive.NewDateTimeFromTime(startDate),
			"$lte": primitive.NewDateTimeFromTime(endDate),
		}},
	}

	var followups = fuMod.GetFollowups(headers.CompanyId, filter, nil)

	response := utils.NewResponseStruct(followups, len(followups))
	response.ToJson(res)
}
func GetFollowUpForContactPerson(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	id := req.URL.Query().Get("id")

	personFollowups := fuMod.GetFollowUpHistoryByContactPerson(headers.CompanyId, id)

	response := utils.NewResponseStruct(personFollowups, len(personFollowups))
	response.ToJson(res)
}

func GetFollowUpForInCharge(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	id, err := strconv.ParseUint(req.URL.Query().Get("id"), 10, 64)

	personFollowups := fuMod.GetFollowUpHistoryByPersonInCharge(headers.CompanyId, id)

	response := utils.NewResponseStruct(personFollowups, len(personFollowups))
	response.ToJson(res)
}

func GetFollowUpsForBill(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	id := req.URL.Query().Get("id")

	personFollowups := fuMod.GetFollowUpHistoryByBill(headers.CompanyId, id)

	response := utils.NewResponseStruct(personFollowups, len(personFollowups))
	response.ToJson(res)
}

func GetFollowUpHistory(res http.ResponseWriter, req *http.Request) {
	headers, err := utils.ResolveHeaders(&req.Header)
	if headers.HandleErrorOrIllegalValues(res, &err) {
		return
	}

	id := req.URL.Query().Get("id")

	personFollowups := fuMod.GetFollowUpHistoryById(headers.CompanyId, id)

	response := utils.NewResponseStruct(personFollowups, len(personFollowups))
	response.ToJson(res)
}
