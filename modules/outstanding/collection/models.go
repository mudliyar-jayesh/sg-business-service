package collection

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	BillNumber    string                  `bson:"Name"`
	PartyName     string                  `bson:"LedgerName"`
	ParentGroup   *string                 `bson:"LedgerGroupName"`
	PendingAmount *models.FloatFromString `bson:"PendingAmount"`
	OpeningAmount *models.FloatFromString `bson:"OpeningAmount"`
	BillDate      *primitive.DateTime     `bson:"BillDate"`
	DueDate       *primitive.DateTime     `bson:"DueDate"`
	DelayDays     *int32
}
