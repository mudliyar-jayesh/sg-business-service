package overview

import (
	"sg-business-service/models"
	"time"
)

type DueTypeEnum string

const (
	NoDue   DueTypeEnum = "noDue"
	Due     DueTypeEnum = "due"
	OverDue DueTypeEnum = "overDue"
)

type OverviewFilter struct {
	Filter               models.RequestFilter
	DeductAdvancePayment bool
	IsDebit              bool
	DueType              DueTypeEnum
	Groups               []string
	Parties              []string
}

type Bill struct {
	BillNumber      string                  `bson:"Name"`
	LedgerName      string                  `bson:"LedgerName"`
	LedgerGroupName *string                 `bson:"LedgerGroupName"`
	OpeningBalance  *models.FloatFromString `bson:"OpeningBalance"`
	ClosingBalance  *models.FloatFromString `bson:"ClosingBalance"`
	BillDate        *time.Time              `bson:"BillDate"`
	DueDate         *time.Time              `bson:"DueDate"`
	IsAdvance       *bool                   `bson:"IsAdvance"`
}

type OutstandingOverview struct {
	PartyName          string
	LedgerGroup        string
	CreditLimit        *float64
	CreditDays         *string
	TotalBills         int
	BillDate           *time.Time
	DueDate            *time.Time
	DelayDays          *uint16
	OpeningAmount      float64
	ClosingAmount      float64
	DueAmount          float64
	OverDueAmount      float64
	ReceivedPercentage *float64
	PendingPercentage  *float64
}
