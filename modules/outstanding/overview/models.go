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
	CreditLimit        *string
	CreditDays         *string
	TotalBills         int
	BillNumber         *string
	BillDate           *time.Time
	DueDate            *time.Time
	DelayDays          *uint16
	OpeningAmount      float64
	ClosingAmount      float64
	DueAmount          float64
	OverDueAmount      float64
	ReceivedPercentage *float64
	PendingPercentage  *float64
	IsAdvance          *bool
	Bills              *[]OutstandingOverview
}

type AgingOverview struct {
	PartyName          string
	LedgerGroup        string
	CreditLimit        *string
	CreditDays         *string
	TotalBills         int
	BillNumber         *string
	BillDate           *time.Time
	DueDate            *time.Time
	DelayDays          *uint16
	OpeningAmount      float64
	ClosingAmount      float64
	Above30          float64
	Above60          float64
	Above90          float64
	Above120          float64
	IsAdvance          *bool
	Bills              *[]AgingOverview
}
