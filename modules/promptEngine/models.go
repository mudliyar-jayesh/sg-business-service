package promptEngine

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sg-business-service/models"
	"time"
)

type ActionableStatus uint8

const (
	Pending ActionableStatus = iota
	Done
	Cancelled
)

type PromptRequest struct {
	StartDateStr string
	EndDateStr   string
	Filter       models.RequestFilter
}

type DecisionRequest struct {
	ActionCode      string
	ApplyOnAllBills bool
	PartyName       *string
	BillNumber      *string
	AmountStr       *string
}

type ActionableOverview struct {
	AssignedToName string
	CreatedByName  string
	StatusName     string
	Task           Actionable
}
type Actionable struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Guid        string             `bson:"guid"`
	CompanyId   string             `bson:"company_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	AssignedTo  int64              `bson:"assigned_to"`
	CreatedBy   int64              `bson:"created_by"`
	CreatedOn   time.Time          `bson:"created_on"`
	Status      ActionableStatus   `bson:"status"`
}

type Action struct {
	Title string
	Code  string
}

type UserPrompt struct {
	Message        string
	Suggestion     string
	SummaryProfile string
	Actions        []Action
	PartyName      *string
	BillNumber     *string
	AmountStr      *string
}
