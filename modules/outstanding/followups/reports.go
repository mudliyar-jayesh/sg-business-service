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

// UpdateFollowUp updates an existing follow-up entry.
// It checks if all associated bills are resolved and updates the follow-up status accordingly.
func UpdateFollowUp(fup FollowUp) error {
	allResolved := true
	for _, bill := range fup.FollowUpBills {
		if bill.Status != Completed {
			allResolved = false
			break
		}
	}
	if allResolved {
		fup.Status = Completed
	}

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

	guid, err := insertFollowUpToDB(fup)
	if err != nil {
		return fmt.Errorf("error inserting follow-up to DB: %w", err)
	}

	fmt.Printf("Follow-up with ID %s inserted to DB\n", guid)
	return nil
}