package followups

import "go.mongodb.org/mongo-driver/bson/primitive"

type ContactPerson struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty"`
	PersonId string				 `bson:"PersonId`
	CompanyId string			 `bson:"CompanyId"`
	Name      string			 `bson:"Name"`
	PartyName string			 `bson:"PartyName"`
	Email     string			 `bson:"Email"`
	PhoneNo   string			 `bson:"PhoneNo"`
}

type FollowUp struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty"`
	RefPrevFollowUpId *string	 `bson:"RefPrevFollowUpId"`
	FollowUpId string   		     `bson:"FollowUpId"`
	ContactPersonId  string 	 `bson:ContactPersonId`
	PersonInChargeId uint32      `bson:PersonInChargeId`
	PartyName string			 `bson:"PartyName"`
	Description string 			 `bson:Description`
	Status FollowUpStatus        `bson:Status`
	FollowUpBills []FollowUpBill `bson:FollowUpBills`
}

type FollowUpStatus int

const (
	Pending FollowUpStatus = iota
	Completed
)

type FollowUpBill struct {
	BillId string `bson:"BillId"`
	Resolved bool `bson:"Resolved"`
}


// ---------- REQUEST MODELS -----------
type FollowUpCreationRequest struct {
	Followup FollowUp
	PointOfContact ContactPerson
}