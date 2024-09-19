package summary

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActionHistory struct {
	Action  string             `bson:"action"`
	Outcome string             `bson:"outcome"`
	Date    primitive.DateTime `bson:"date"`
}

type OutstandingSummary struct {
	ClientId          string          `bson:"client_id"`
	LedgerName        string          `bson:"ledger_name"`
	TotalTransactions int             `bson:"total_transactions"`
	TotalDelayed      int             `bson:"total_delayed"`
	DelayPercentage   float64         `bson:"delay_percentage"`
	AmountDue         float64         `bson:"amount_due"`
	AverageDelayDays  int             `bson:"average_delay_days"`
	LastAction        string          `bson:"last_action"`
	LastOutcome       string          `bson:"last_outcome"`
	ActionHistory     []ActionHistory `bson:"action_history"`
}
