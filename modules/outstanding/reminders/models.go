package reminders

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	osMod "sg-business-service/modules/outstanding"
)

type OutstandingTemplate struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CompanyId    string             `bson:"company_id"`
	TemplateName string             `bson:"template_name"`
	HtmlContent  string             `bson:"html_content"`
}

type ReminderBody struct {
	PartyName      string
	Address        string
	TotalAmount    float64
	TotalAmountStr string
	Bills          []osMod.Bill
}
