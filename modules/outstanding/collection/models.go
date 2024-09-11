package collection

import (
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"sg-business-service/models"
)

type CollectionFilter struct {
	StartDateStr string
	EndDateStr   string
	Parties      []string
	Groups       []string
	Filter       models.RequestFilter
}

type CollectionOverview struct {
	BillNumber  string  `bson:"Name"`
	PartyName   string  `bson:"LedgerName"`
	ParentGroup *string `bson:"LedgerGroupName"`
	/*DueDate       *primitive.DateTime `bson:"DueDate"`
	BillDate      *primitive.DateTime `bson:"BillDate"`
	DelayDays     *int32   */
	OpeningAmount *float64 `bson:"OpeningBal"`
	PendingAmount *float64 `bson:"ClosingBal"`
}
