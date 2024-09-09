package reminders

import (
	osMod "sg-business-service/modules/outstanding"
)

type ReminderBody struct {
    PartyName string
    Address string
    TotalAmount float64
    Bills []osMod.Bill
}



