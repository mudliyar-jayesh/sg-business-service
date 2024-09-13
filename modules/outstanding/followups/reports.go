package followups

import (
	"fmt"
	"time"
)

func GetContactPersons(companyId string, partyName string) ([] ContactPerson){
	contactPersons := getContactPersonList(companyId, partyName)
	return contactPersons;
}

func GetFollowUpList(companyId string, partyName string) []FollowUp{
	return getFollowupListByParty(companyId, partyName)	
}

func UpdateFollowUp(fup FollowUp) error {
	/***
		Updates a followup entry
	*/

	// First check if the followup exists

	// If all the bills are resolved then set the followup status to resolved.

	// Update metadata
	currentDt := time.Now()
	fup.LastUpdated = &currentDt 

	// Return the updated followup to user
	_, err := insertFollowUpToDB(fup)	
	return err
}

func CreateFollowUp(fup FollowUp, cperson *ContactPerson) error {
	/***
		Creates a follow up and creates contact person is not an existing provided

		Below are the cases that are handled by this function: 
		 case 1: len(FollowUp.ContactPersonId) < 1 and cperson is not null then we have to create a new contact person.
	 	 case 2: len(FollowUp.ContactPersonId) > 1 then it means then it means an already registered contact person is handling the followup.
	*/

	// Register contact person if not already
	if len(fup.ContactPersonId) <= 1 && cperson != nil {

		cpersonId, err := createContactPerson(*cperson)

		if err != nil {
			fmt.Printf("Error occured while creating a new followup person")
		}
		fup.ContactPersonId = cpersonId
	}

	// before inserting any followup we check what are the previous followups for this

	// Inject datetime metadata before sending to DB
	created := time.Now()
	fup.Created = &created 

	// At final stage insert the  followup	
	guid, err := insertFollowUpToDB(fup)

	if err != nil {
		return err
	}

	fmt.Printf("Followup with %s inserted to DB", guid)
	return nil;
}