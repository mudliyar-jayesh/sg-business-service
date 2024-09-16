package followups

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"sg-business-service/models"
)

type ContactPerson struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty"`
	PersonId  string              `bson:"PersonId"`
	CompanyId string              `bson:"CompanyId"`
	Name      string              `bson:"Name"`
	PartyName string              `bson:"PartyName"`
	Email     string              `bson:"Email"`
	PhoneNo   string              `bson:"PhoneNo"`
}

type FollowUp struct {
	ID                *primitive.ObjectID `bson:"_id,omitempty"`
	Created           *time.Time          `bson:"CreateDate"`
	LastUpdated       *time.Time          `bson:"LastUpdated"`
	CompanyId         string              `bson:"CompanyId"`
	RefPrevFollowUpId *string             `bson:"RefPrevFollowUpId"`
	FollowUpId        string              `bson:"FollowUpId"`
	ContactPersonId   string              `bson:"ContactPersonId"`
	PersonInChargeId  uint64              `bson:"PersonInChargeId"`
	PartyName         string              `bson:"PartyName"`
	Description       string              `bson:"Description"`
	Status            FollowUpStatus      `bson:"Status"`
	FollowUpBills     []FollowUpBill      `bson:"FollowUpBills"`
	NextFollowUpDate  *time.Time          `bson:"NextFollowUpDate"`
}

type FollowUpStatus int

const (
	Pending FollowUpStatus = iota
	Scheduled
	Completed
)

func GetFollowUpStatusMappings() map[string]int {
	return map[string]int{
		"Pending":   0,
		"Scheduled": 1,
		"Completed": 2,
	}
}

type FollowUpFilter struct {
	StartDateStr string
	EndDateStr   string
	Filter       models.RequestFilter
}
type FollowUpHistory struct {
	PartyName        string
	CreationDate     time.Time
	UpdationDate     time.Time
	NextFollowUpDate *time.Time
	PersonInCharge   string
	PersonInChargeId uint64
	PocName          string
	PocEmail         string
	PocMobile        string
	Description      string
	FollowUpBills    []FollowUpBill
	TotalCount       int32
	PendingCount     int32
	ScheduledCount   int32
	CompleteCount    int32
}

type FollowUpOverview struct {
	Name           string
	Amount         float64
	TotalCount     int32
	PendingCount   int32
	ScheduledCount int32
	CompleteCount  int32
}
type FollowUpBill struct {
	BillId string         `bson:"BillId"`
	Status FollowUpStatus `bson:"Status"`
}

// ---------- REQUEST MODELS -----------
type FollowUpCreationRequest struct {
	Followup       FollowUp
	PointOfContact *ContactPerson
}
