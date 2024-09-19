package summary

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sg-business-service/models"
	osMod "sg-business-service/modules/outstanding"
	"sg-business-service/utils"
	"time"
)

// Function to fetch MetaBills and group by LedgerName
func CalculateOutstandingSummary(companyId string) {

	var requestFilter = models.RequestFilter{}

	var bills = osMod.GetBills(companyId, requestFilter, true, nil)

	groupedBills := utils.GroupFor(bills, func(bill osMod.MetaBill) string {
		return bill.PartyName
	})

	var summaries []OutstandingSummary
	for ledgerName, bills := range groupedBills {
		summary := populateOutstandingSummary(companyId, ledgerName, bills)
		summaries = append(summaries, summary)
	}

	var convertedSummaries = ConvertToInterfaceSlice(summaries)
	var handler = GetCollection()
	handler.InsertMany(convertedSummaries)
}

// Function to populate OutstandingSummary from grouped MetaBill entries
func populateOutstandingSummary(companyId, ledgerName string, bills []osMod.MetaBill) OutstandingSummary {
	totalTransactions := len(bills)
	totalDelayed := 0
	totalAmountDue := 0.0
	totalDelayDays := 0
	actionHistory := []ActionHistory{}

	currentTime := primitive.NewDateTimeFromTime(time.Now())

	for _, bill := range bills {
		if bill.PendingAmount != nil {
			totalAmountDue += bill.PendingAmount.Value // Assuming FloatFromString.Value contains the float value
		}

		// Calculate delay if DueDate is before the current date
		if bill.DueDate != nil && bill.DueDate.Time().Before(currentTime.Time()) {
			totalDelayed++
			delayDays := int(currentTime.Time().Sub(bill.DueDate.Time()).Hours() / 24)
			totalDelayDays += delayDays
		}

		// Simulating ActionHistory TODO: Remove later
		actionHistory = append(actionHistory, ActionHistory{
			Action:  "Check Payment Status", // Example action
			Outcome: "Pending",              // Example outcome
			Date:    primitive.NewDateTimeFromTime(bill.BillDate.Time()),
		})
	}

	// Calculate delay percentage and average delay days
	delayPercentage := 0.0
	averageDelayDays := 0

	if totalTransactions > 0 {
		delayPercentage = (float64(totalDelayed) / float64(totalTransactions)) * 100
	}

	if totalDelayed > 0 {
		averageDelayDays = totalDelayDays / totalDelayed
	}

	// Example of determining last action and outcome (simplified)
	lastAction := "Send Reminder"
	lastOutcome := "Payment Pending"

	return OutstandingSummary{
		ClientId:          companyId,
		LedgerName:        ledgerName,
		TotalTransactions: totalTransactions,
		TotalDelayed:      totalDelayed,
		DelayPercentage:   delayPercentage,
		AmountDue:         totalAmountDue,
		AverageDelayDays:  averageDelayDays,
		LastAction:        lastAction,
		LastOutcome:       lastOutcome,
		ActionHistory:     actionHistory,
	}
}
