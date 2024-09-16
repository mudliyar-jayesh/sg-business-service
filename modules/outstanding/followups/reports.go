// Package followups provides functionality for managing follow-ups and contact persons.
package followups

import (
	"fmt"
	"time"
)

// GetContactPersons retrieves a list of contact persons for a given company and party.
func GetContactPersons(companyID, partyName string) []ContactPerson {
	return getContactPersonList(companyID, partyName)
}

// GetFollowUpList retrieves a list of follow-ups for a given company and party.
func GetFollowUpList(companyID, partyName string) []FollowUp {
	return getFollowupListByParty(companyID, partyName)
}

// GetFollowUpHistoryById retrieves all the connected followups
func GetFollowUpHistoryById(companyID, followUpId string) []FollowUp {
	return getFollowupHistoryById(companyID, followUpId)
}

// GetFollowUpHistoryByBill retrieves all the followups for a party bill
func GetFollowUpHistoryByBill(companyID, billId string) []FollowUp {
	return getFollowupHistoryByBillId(companyID, billId)
}

// GetFollowUpHistoryByContactPerson retrieves all the followups for a contact person
func GetFollowUpHistoryByContactPerson(companyID, contactPersonId string) []FollowUp {
	return getFollowupHistoryByContactPerson(companyID, contactPersonId)
}

// GetFollowUpHistoryByPersonInCharge retrieves all the followups attented by a person in charge
func GetFollowUpHistoryByPersonInCharge(companyID string, personInChargeId uint64) []FollowUp {
	return getFollowUpHistoryByPersonInCharge(companyID, personInChargeId)
}

// calibrateFollowupStatus: modifies the status of a Followup object based on the status of it's bills
//
//	it will assign the status which is found in majority of the bills.
func calibrateFollowupStatus(fup *FollowUp) {
	statusCountMap := make(map[FollowUpStatus]int)

	for _, bill := range fup.FollowUpBills {
		if count, exists := statusCountMap[bill.Status]; exists {
			statusCountMap[bill.Status] = count + 1
		} else {
			statusCountMap[bill.Status] = 1
		}
	}

	var maxCountedStatus FollowUpStatus

	maxCount := 0

	for status, count := range statusCountMap {
		if count > maxCount {
			maxCount = count
			maxCountedStatus = status
		}
	}

	fup.Status = maxCountedStatus
}

// UpdateFollowUp updates an existing follow-up entry.
// It checks if all associated bills are resolved and updates the follow-up status accordingly.
func UpdateFollowUp(fup FollowUp) error {

	calibrateFollowupStatus(&fup)

	currentDt := time.Now()
	fup.LastUpdated = &currentDt

	return updateFollowup(fup)
}

// CreateFollowUp creates a new follow-up entry and optionally creates a new contact person.
// It handles two cases:
// 1. If ContactPersonId is empty and cperson is not nil, it creates a new contact person.
// 2. If ContactPersonId is not empty, it uses the existing contact person.
func CreateFollowUp(fup FollowUp, cperson *ContactPerson) error {
	if fup.ContactPersonId == "" && cperson != nil {
		cpersonID, err := createContactPerson(*cperson)
		if err != nil {
			return fmt.Errorf("error creating new contact person: %w", err)
		}
		fup.ContactPersonId = cpersonID
	}

	created := time.Now()
	fup.Created = &created
	fup.LastUpdated = &created

	calibrateFollowupStatus(&fup)

	guid, err := insertFollowUpToDB(fup)
	if err != nil {
		return fmt.Errorf("error inserting follow-up to DB: %w", err)
	}
	fmt.Printf("Follow-up with ID %s inserted to DB\n", guid)
	return nil
}
