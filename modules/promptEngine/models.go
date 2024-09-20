package promptEngine

import (
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

type Actionable struct {
	Title       string
	Description string
	CreatedBy   uint64
	CreatedOn   time.Time
	Status      ActionableStatus
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
}
