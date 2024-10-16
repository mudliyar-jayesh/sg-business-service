package prompts

import (
	"fmt"
	"github.com/leekchan/accounting"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sg-business-service/models"
	osMod "sg-business-service/modules/outstanding"
	osSummary "sg-business-service/modules/outstanding/summary"
	"sg-business-service/modules/promptEngine"
	"sg-business-service/utils"
	"time"
)

var ac = accounting.Accounting{Symbol: "₹", Precision: 2}

const (
	SendReminderAction string = "send_reminder"
	TeamFollowUpAction string = "team_follow_up"
)

var CollectionActions = []promptEngine.Action{
	{
		Title: "Send Reminder",
		Code:  SendReminderAction,
	},
	{
		Title: "Notify Team To Follow-Up",
		Code:  TeamFollowUpAction,
	},
	{
		Title: "Ignore",
		Code:  "ignore",
	},
}

func GetCollectionPrompt(companyId, fromDate, toDate string, requestFilter models.RequestFilter, partyWise bool) []promptEngine.UserPrompt {
	startDate, err := time.Parse("02-01-2006", fromDate)
	if err != nil {
		return make([]promptEngine.UserPrompt, 0)
	}
	endDate, err := time.Parse("02-01-2006", toDate)
	if err != nil {
		return make([]promptEngine.UserPrompt, 0)
	}

	dbFilter := []bson.M{
		{"BillCreditPeriod.DueDate": bson.M{"$exists": true, "$ne": nil}},
		{"BillCreditPeriod.DueDate": bson.M{
			"$gte": primitive.NewDateTimeFromTime(startDate),
			"$lte": primitive.NewDateTimeFromTime(endDate),
		}},
	}
	var bills = osMod.GetBills(companyId, requestFilter, true, dbFilter)

	var partyWiseBills = utils.GroupFor(bills, func(bill osMod.MetaBill) string {
		return bill.PartyName
	})

	var prompts []promptEngine.UserPrompt
	for partyName, partyBills := range partyWiseBills {
		var summaries = osSummary.GetSummaryByPartyName(companyId, partyName)
		var summaryStr string = "Summary\n"
		for _, summary := range summaries {
			var message = fmt.Sprintf("\nTotal Transactions: %v\n Average Delay Days: %v\nDelay Percentage: %v\n", summary.TotalTransactions, summary.AverageDelayDays, summary.DelayPercentage)
			summaryStr += message
		}

		var totalAmount float64
		for _, bill := range partyBills {
			// TODO: later add the prediction call and check
			var amount = bill.OpeningAmount.Value
			if bill.PendingAmount != nil {
				amount = bill.PendingAmount.Value
			}
			totalAmount += amount
			var amountStr = ac.FormatMoney(amount)

			var dueDate = bill.BillDate.Time()
			if bill.DueDate != nil {
				dueDate = bill.DueDate.Time()
			}
			today := time.Now().UTC()

			// Calculate the difference
			diff := today.Sub(dueDate)

			// Get the difference in days
			days := int32(diff.Hours() / 24)

			var promptMessage = fmt.Sprintf("%v has pending amount of %v on bill number: %v ", partyName, amountStr, bill.BillNumber)
			if days > 0 {
				var delayMessage = fmt.Sprintf(" with a delay of %v days", days)
				promptMessage += delayMessage
			}
			if !partyWise {
				var userPrompt = promptEngine.UserPrompt{
					Message:        promptMessage,
					SummaryProfile: summaryStr,
					Actions:        CollectionActions,
					PartyName:      &partyName,
					BillNumber:     &bill.BillNumber,
					AmountStr:      &amountStr,
				}
				prompts = append(prompts, userPrompt)
			}
		}
		if partyWise {
			var totalAmountStr = ac.FormatMoney(totalAmount)
			var promptMessage = fmt.Sprintf("%v has pending amount of %v", partyName, totalAmountStr)
			var userPrompt = promptEngine.UserPrompt{
				Message:        promptMessage,
				SummaryProfile: summaryStr,
				Actions:        CollectionActions,
				PartyName:      &partyName,
				AmountStr:      &totalAmountStr,
			}
			prompts = append(prompts, userPrompt)
		}
	}

	return prompts
}
