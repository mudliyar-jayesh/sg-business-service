package overview

import (
	"go.mongodb.org/mongo-driver/bson"
	"sg-business-service/models"
	"sg-business-service/modules/ledgers"
	osMod "sg-business-service/modules/outstanding"
	osSettingMod "sg-business-service/modules/outstanding/settings"
	"sg-business-service/utils"

	"math"
	"sync"
	"time"

	"log"
)

const (
	PartyWise      string = "party-wise"
	GroupWise      string = "group-wise"
	CreditLimit    string = "credit-limit-wise"
	CreditPeriod   string = "credit-period-wise"
	BillWise       string = "bill-wise"
	BillDateWise   string = "bill-date-wise"
	DueDateWise    string = "due-date-wise"
	OpeningWise    string = "opening-wise"
	ClosingWise    string = "closing-wise"
	DueWise        string = "due-wise"
	OverDueWise    string = "over-due-wise"
	TotalBillsWise string = "bill-count-wise"
)

func GetPartyWiseOverview(companyId string, filter OverviewFilter) []OutstandingOverview {

	// Ledger Groups Filter
	var groups = getParentGroups(companyId, filter.IsDebit)
	if len(filter.Groups) > 0 {
		groups = utils.Intersection(groups, filter.Groups)
	}
	filter.Groups = groups

	var laterSortKeys = make([]string, 4)
	laterSortKeys[0] = OpeningWise
	laterSortKeys[1] = ClosingWise
	laterSortKeys[2] = DueWise
	laterSortKeys[3] = OverDueWise

	var usePagination bool = !utils.ContainsString(laterSortKeys, filter.Filter.SortKey)
	log.Printf("\n [!] Use Pagination  :: %v  \n", usePagination)

	var ledgerFilter = models.RequestFilter{}
	ledgerFilter.Batch.Apply = usePagination
	ledgerFilter.Batch.Limit = filter.Filter.Batch.Limit
	ledgerFilter.Batch.Offset = filter.Filter.Batch.Offset

	var sortKey string
	switch filter.Filter.SortKey {
	case GroupWise:
		sortKey = "Group"
	case CreditLimit:
		sortKey = "CreditLimit"
	case CreditPeriod:
		sortKey = "CreditPeriod"
	default:
		sortKey = "Name"
	}

	ledgerFilter.SortKey = sortKey
	ledgerFilter.SortOrder = filter.Filter.SortOrder

	var ledgerAdditionalFilter = []bson.M{
		{
			"Group": bson.M{
				"$in": groups,
			},
		},
	}

	if len(filter.Parties) > 0 {
		ledgerAdditionalFilter = append(ledgerAdditionalFilter, bson.M{
			"Name": bson.M{
				"$in": filter.Parties,
			},
		})
	}

	// Search Filter
	var useLedgerSearch = filter.Filter.SearchKey == "party-wise" || filter.Filter.SearchKey == "group-wise"
	log.Printf("\n [!] Use Ledger Search  :: %v  \n", useLedgerSearch)
	if len(filter.Filter.SearchText) > 0 && useLedgerSearch {
		var searchKey string
		switch filter.Filter.SearchKey {
		case GroupWise:
			searchKey = "Group"
		case CreditLimit:
			searchKey = "CreditLimit"
		case CreditPeriod:
			searchKey = "CreditPeriod"
		default:
			searchKey = "Name"
		}
		var searchField = osMod.GetFieldBySearchKey(searchKey)
		var searchFilter = utils.GenerateSearchFilter(filter.Filter.SearchText, searchField)
		ledgerAdditionalFilter = append(ledgerAdditionalFilter, bson.M{
			"$and": &searchFilter,
		})
		log.Printf("\n [!] Search Added:: %v  \n", useLedgerSearch)
	}

	var parties = ledgers.GetLedgers(companyId, ledgerFilter, ledgerAdditionalFilter)
	log.Printf("\n [+] Parties Found :: %v  \n", len(parties))

	var partyByName = utils.ToDict(parties, func(party ledgers.MetaLedger) string {
		return party.Name
	})

	settings, settingsErr := osSettingMod.GetAllSettings(companyId)
	if settingsErr != nil {
		return make([]OutstandingOverview, 0)
	}

	setting := settings[0]
	istLocation, _ := time.LoadLocation("Asia/Kolkata")

	var partySummary = make(map[string][]OutstandingOverview, 0)

	var mutex sync.Mutex

	var batchFunc = func(partyChunk []ledgers.MetaLedger, wait *sync.WaitGroup) {
		defer wait.Done()
		var partyNames = utils.Select(partyChunk, func(party ledgers.MetaLedger) string {
			return party.Name
		})
		log.Printf("\n [+] Parties Names :: %v  \n", len(partyNames))

		billFilter := []bson.M{
			{
				"LedgerName": bson.M{
					"$in": partyNames,
				},
			},
		}

		var billDbFilter = filter
		billDbFilter.Filter.SortKey = "Name"
		var bills = getBills(companyId, billDbFilter, &billFilter)
		log.Printf("\n [+] Bills  :: %v  \n", len(bills))

		if len(bills) < 1 {
			mutex.Lock()
			for _, partyName := range partyNames {
				partyInfo := partyByName[partyName]
				var overview = OutstandingOverview{
					PartyName:     partyName,
					LedgerGroup:   partyInfo.Group,
					CreditDays:    partyInfo.CreditPeriod,
					OpeningAmount: 0,
					ClosingAmount: 0,
					DueAmount:     0,
					OverDueAmount: 0,
				}
				partySummary[partyName] = make([]OutstandingOverview, 1)
				partySummary[partyName][0] = overview

			}
			mutex.Unlock()
		}
		today := time.Now().UTC()
		for _, bill := range bills {
			var overview = OutstandingOverview{
				PartyName:     bill.LedgerName,
				LedgerGroup:   *bill.LedgerGroupName,
				OpeningAmount: bill.OpeningBalance.Value,
				ClosingAmount: bill.ClosingBalance.Value,
				DueAmount:     0,
				OverDueAmount: 0,
			}
			billDate := bill.BillDate.In(istLocation)
			dueDate := billDate
			if bill.DueDate != nil {
				dueDate = *bill.DueDate
			}
			diff := today.Sub(dueDate)

			// Get the difference in days
			days := uint16(diff.Hours() / 24)

			overview.BillDate = &billDate
			overview.DueDate = &dueDate
			overview.DelayDays = &days
			overview.IsAdvance = bill.IsAdvance

			if days >= 0 && days <= uint16(setting.OverDueDays) {
				overview.DueAmount = overview.ClosingAmount
			} else if days > uint16(setting.OverDueDays) {
				overview.OverDueAmount = overview.ClosingAmount
			}
			mutex.Lock()
			_, exists := partySummary[bill.LedgerName]
			if !exists {
				partySummary[bill.LedgerName] = make([]OutstandingOverview, 0)
			}

			partySummary[bill.LedgerName] = append(partySummary[bill.LedgerName], overview)
			mutex.Unlock()
		}
	}

	var totalParites = float64(len(parties))
	var batchSize = math.Ceil(totalParites * 0.25)
	var batchCount = math.Ceil(totalParites / batchSize)

	var wg sync.WaitGroup

	for batchNumber := 0; batchNumber < int(batchCount); batchNumber++ {
		wg.Add(1)
		var startIndex = int(batchSize) * batchNumber
		var endIndex = startIndex + int(batchSize)
		var partyChunk = parties[startIndex:endIndex]
		go batchFunc(partyChunk, &wg)
	}
	wg.Wait()

	log.Printf("\n [+] Parties Summarized :: %v  \n", len(partySummary))
	var outstandingOverview []OutstandingOverview
	for partyName, bills := range partySummary {
		party, exists := partyByName[partyName]
		if !exists {
			continue
		}

		var overview = OutstandingOverview{
			PartyName:          party.Name,
			LedgerGroup:        party.Group,
			CreditLimit:        party.CreditLimit,
			CreditDays:         party.CreditPeriod,
			OpeningAmount:      0,
			ClosingAmount:      0,
			DueAmount:          0,
			OverDueAmount:      0,
			PendingPercentage:  nil,
			ReceivedPercentage: nil,
			TotalBills:         len(bills),
		}

		for _, bill := range bills {
			overview.OpeningAmount += bill.OpeningAmount
			if filter.DeductAdvancePayment && bill.IsAdvance != nil && *bill.IsAdvance {
				overview.ClosingAmount -= bill.ClosingAmount
				overview.DueAmount -= bill.DueAmount
				overview.OverDueAmount -= bill.OverDueAmount
				continue
			}

			overview.ClosingAmount += bill.ClosingAmount
			overview.DueAmount += bill.DueAmount
			overview.OverDueAmount += bill.OverDueAmount
		}
		var receivedAmount = overview.OpeningAmount - overview.ClosingAmount
		var percentReceived float64 = 0
		if overview.OpeningAmount > 0 {
			percentReceived = (receivedAmount / overview.OpeningAmount) * 100
		}
		var percentPending = 100 - percentReceived
		overview.PendingPercentage = &percentPending
		overview.ReceivedPercentage = &percentReceived
		outstandingOverview = append(outstandingOverview, overview)
	}
	if !usePagination {
		sortAsc := utils.GetValueBySortOrder(filter.Filter.SortOrder) == 1
		switch filter.Filter.SortKey {
		case OpeningWise:
			utils.SortByField(outstandingOverview, "OpeningAmount", sortAsc)
		case ClosingWise:
			utils.SortByField(outstandingOverview, "ClosingAmount", sortAsc)
		case DueWise:
			utils.SortByField(outstandingOverview, "DueAmount", sortAsc)
		case OverDueWise:
			utils.SortByField(outstandingOverview, "OverDueAmount", sortAsc)
		case DueDateWise:
			utils.SortByField(outstandingOverview, "DueDate", sortAsc)
		case BillDateWise:
			utils.SortByField(outstandingOverview, "BillDate", sortAsc)
		case TotalBillsWise:
			utils.SortByField(outstandingOverview, "TotalBills", sortAsc)
		default:
			utils.SortByField(outstandingOverview, "PartyName", sortAsc)
		}

		outstandingOverview = utils.Paginate(outstandingOverview, int(filter.Filter.Batch.Limit), int(filter.Filter.Batch.Offset))
	}
	log.Printf("\n [+] Party Overviews :: %v  \n", len(outstandingOverview))
	log.Println("-------------------------------------------------------")
	// sorting by keys
	return outstandingOverview
}

func BillWiseOverview(filter OverviewFilter) {

}
