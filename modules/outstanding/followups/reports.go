package followups

import (
	"fmt"
)

func createFollowUp(fup FollowUp, cperson *ContactPerson) {
	/***
		Creates a follow up and creates contact person is not an existing provided

		Below are the cases that are handled by this function: 
		 case 1: len(FollowUp.ContactPersonId) < 1 and cperson is not null then we have to create a new contact person.
	 	 case 2: len(FollowUp.ContactPersonId) > 1 then it means then it means an already registered contact person is handling the followup.
	*/

	// before inserting any followup we check what are the previous followups for this

	// At final stage insert the  followup	
	guid, err := insertFollowUpToDB(fup)

	if err != nil {
		fmt.Println(err);
	}

	fmt.Printf("Followup with %s inserted to DB", guid)
}