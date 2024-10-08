package overview

import (
  "sg-business-service/models"
  "sg-business-service/modules/ledgers"
  osMod "sg-business-service/modules/outstanding"
  osSettingMod "sg-business-service/modules/outstanding/settings"
  "sg-business-service/utils"

  "go.mongodb.org/mongo-driver/bson"

  "math"
  "sync"
  "time"
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
  DelayWise      string = "delay-days-wise"

  Above30Wise      string = "above-30-wise"
  Above60Wise      string = "above-60-wise"
  Above90Wise      string = "above-90-wise"
  Above120Wise      string = "above-120-wise"
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
        "$in": filter.Groups,
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
  }

  var parties = ledgers.GetLedgers(companyId, ledgerFilter, ledgerAdditionalFilter)

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

    if len(bills) < 1 {
      mutex.Lock()
      for _, partyName := range partyNames {
        partyInfo := partyByName[partyName]
        var overview = OutstandingOverview{
          PartyName:     partyName,
          LedgerGroup:   partyInfo.Group,
          CreditDays:    partyInfo.CreditPeriod,
          CreditLimit:   partyInfo.CreditLimit,
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
        BillNumber:    &bill.BillNumber,
        LedgerGroup:   *bill.LedgerGroupName,
        OpeningAmount: bill.OpeningBalance.Value,
        ClosingAmount: bill.ClosingBalance.Value,
        DueAmount:     0,
        OverDueAmount: 0,
      }
      billDate := bill.BillDate.In(istLocation)
      dueDate := billDate
      if bill.DueDate != nil {
        dueDate = (*bill.DueDate).In(istLocation)
      }
      diff := today.Sub(dueDate)

      // Get the difference in days
      days := int(diff.Hours() / 24)

      // Ensure DelayDays is never negative. If days < 0, set it to 0 (or handle according to your logic).
      var delayDays uint16
      if days > 0 {
        delayDays = uint16(days)
      } else {
        delayDays = 0
      }
      overview.BillDate = &billDate
      overview.DueDate = &dueDate
      overview.DelayDays = &delayDays
      overview.IsAdvance = bill.IsAdvance

      if delayDays >= 0 && delayDays <= uint16(setting.OverDueDays) {
        overview.DueAmount = overview.ClosingAmount
      } else if delayDays > uint16(setting.OverDueDays) {
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
    if endIndex >= len(parties) {
      endIndex = len(parties) - 1
    }
    var partyChunk = parties[startIndex:endIndex]
    go batchFunc(partyChunk, &wg)
  }
  wg.Wait()

  var outstandingOverview []OutstandingOverview
  for _, partyLedger := range parties {
    var overview = OutstandingOverview{
      PartyName:          partyLedger.Name,
      LedgerGroup:        partyLedger.Group,
      CreditLimit:        partyLedger.CreditLimit,
      CreditDays:         partyLedger.CreditPeriod,
      OpeningAmount:      0,
      ClosingAmount:      0,
      DueAmount:          0,
      OverDueAmount:      0,
      PendingPercentage:  nil,
      ReceivedPercentage: nil,
      TotalBills:         0,
    }
    summary, exist := partySummary[partyLedger.Name]
    if exist {
      overview.Bills = &summary
      var billCount int
      for _, bill := range summary {
        overview.OpeningAmount += bill.OpeningAmount
        billCount += 1
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
      overview.TotalBills = billCount
    }

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
  // sorting by keys
  return outstandingOverview
}

func GetBillWiseOverview(companyId string, filter OverviewFilter) []OutstandingOverview {
  // Ledger Groups Filter
  var groups = getParentGroups(companyId, filter.IsDebit)
  if len(filter.Groups) > 0 {
    groups = utils.Intersection(groups, filter.Groups)
  }
  filter.Groups = groups

  var laterSortKeys = make([]string, 6)
  laterSortKeys[0] = DueWise
  laterSortKeys[1] = DelayWise
  laterSortKeys[2] = OverDueWise
  laterSortKeys[3] = DueDateWise
  laterSortKeys[4] = CreditLimit
  laterSortKeys[5] = CreditPeriod

  var usePagination bool = !utils.ContainsString(laterSortKeys, filter.Filter.SortKey)

  var sortKey string
  switch filter.Filter.SortKey {
  case GroupWise:
    sortKey = "LedgerGroupName"
  case BillDateWise:
    sortKey = "BillDate.Date"
  case CreditPeriod:
    sortKey = "CreditPeriod"
  default:
    sortKey = "LedgerName"
  }

  var billDbFilter = filter
  billDbFilter.Groups = filter.Groups
  billDbFilter.Filter.Batch.Apply = usePagination
  billDbFilter.Filter.Batch.Limit = filter.Filter.Batch.Limit
  billDbFilter.Filter.Batch.Offset = filter.Filter.Batch.Offset
  billDbFilter.Filter.SortKey = sortKey
  billDbFilter.Filter.SortOrder = filter.Filter.SortOrder
  billDbFilter.Parties = filter.Parties

  var billAdditionalFilter = []bson.M{
    {
      "LedgerGroupName": bson.M{
        "$in": billDbFilter.Groups,
      },
    },
  }

  if len(billDbFilter.Parties) > 0 {
    billAdditionalFilter = append(billAdditionalFilter, bson.M{
      "LedgerName": bson.M{
        "$in": billDbFilter.Parties,
      },
    })
  }
  if len(filter.Filter.SearchText) > 0 {
    var searchKey string
    switch filter.Filter.SearchKey {
    case GroupWise:
      searchKey = "LedgerGroupName"
    case BillWise:
      searchKey = "Name"
    case BillDateWise:
      searchKey = "BillDate.Date"
    default:
      searchKey = "Name"
    }
    var searchField = osMod.GetFieldBySearchKey(searchKey)
    var searchFilter = utils.GenerateSearchFilter(filter.Filter.SearchText, searchField)
    billAdditionalFilter = append(billAdditionalFilter, bson.M{
      "$and": &searchFilter,
    })
  }

  settings, settingsErr := osSettingMod.GetAllSettings(companyId)
  if settingsErr != nil {
    return make([]OutstandingOverview, 0)
  }
  setting := settings[0]
  istLocation, _ := time.LoadLocation("Asia/Kolkata")

  var bills = getBills(companyId, billDbFilter, &billAdditionalFilter)

  var partyNames []string
  var billSummary []OutstandingOverview
  today := time.Now().UTC()
  for _, bill := range bills {
    var overview = OutstandingOverview{
      PartyName:     bill.LedgerName,
      LedgerGroup:   *bill.LedgerGroupName,
      OpeningAmount: bill.OpeningBalance.Value,
      ClosingAmount: bill.ClosingBalance.Value,
      BillNumber:    &bill.BillNumber,
      DueAmount:     0,
      OverDueAmount: 0,
    }
    partyNames = append(partyNames, bill.LedgerName)

    billDate := bill.BillDate.In(istLocation)
    dueDate := billDate
    if bill.DueDate != nil {
      dueDate = (*bill.DueDate).In(istLocation)
    }
    diff := today.Sub(dueDate)

    // Get the difference in days
    days := int(diff.Hours() / 24)

    // Ensure DelayDays is never negative. If days < 0, set it to 0 (or handle according to your logic).
    var delayDays uint16
    if days > 0 {
      delayDays = uint16(days)
    } else {
      delayDays = 0
    }

    overview.BillDate = &billDate
    overview.DueDate = &dueDate
    overview.DelayDays = &delayDays
    overview.IsAdvance = bill.IsAdvance

    if delayDays >= 0 && delayDays <= uint16(setting.OverDueDays) {
      overview.DueAmount = overview.ClosingAmount
    } else if delayDays > uint16(setting.OverDueDays) {
      overview.OverDueAmount = overview.ClosingAmount
    }

    billSummary = append(billSummary, overview)
  }

  var distinctPartyNames = utils.Distinct(partyNames)

  var parties = ledgers.GetByNames(companyId, distinctPartyNames)
  var partyByName = utils.ToDict(parties, func(party ledgers.MetaLedger) string {
    return party.Name
  })

  var billOverview []OutstandingOverview
  for _, summary := range billSummary {
    party, exists := partyByName[summary.PartyName]
    if !exists {
      continue
    }

    var overview = summary
    overview.CreditLimit = party.CreditLimit
    overview.CreditDays = party.CreditPeriod

    billOverview = append(billOverview, overview)
  }

  if !usePagination {
    sortAsc := utils.GetValueBySortOrder(filter.Filter.SortOrder) == 1
    switch filter.Filter.SortKey {
    case DueWise:
      utils.SortByField(billOverview, "DueAmount", sortAsc)
    case DelayWise:
      utils.SortByField(billOverview, "DelayDays", sortAsc)
    case DueDateWise:
      utils.SortByField(billOverview, "DueDate", sortAsc) 
    case OverDueWise:
      utils.SortByField(billOverview, "OverDueAmount", sortAsc)
    case CreditLimit:
      utils.SortByField(billOverview, "CreditLimit", sortAsc)
    case CreditPeriod:
      utils.SortByField(billOverview, "CreditDays", sortAsc)
    default:
      utils.SortByField(billOverview, "PartyName", sortAsc)
    }

    billOverview = utils.Paginate(billOverview, int(filter.Filter.Batch.Limit), int(filter.Filter.Batch.Offset))
  }

  return billOverview
}
func GetAgingOverview(companyId string, useAgingRange bool, filter OverviewFilter) []AgingOverview {
  // Ledger Groups Filter
  var groups = getParentGroups(companyId, filter.IsDebit)
  if len(filter.Groups) > 0 {
    groups = utils.Intersection(groups, filter.Groups)
  }
  filter.Groups = groups

  var laterSortKeys = []string{Above30Wise, Above60Wise, Above90Wise, Above120Wise}

  var usePagination = !utils.ContainsString(laterSortKeys, filter.Filter.SortKey)
  var ledgerFilter = models.RequestFilter{
    Batch: models.Pagination{
      Apply:  usePagination,
      Limit:  filter.Filter.Batch.Limit,
      Offset: filter.Filter.Batch.Offset,
    },
    SortKey:   "Name",
    SortOrder: filter.Filter.SortOrder,
  }

  var ledgerAdditionalFilter = []bson.M{
    {"Group": bson.M{"$in": filter.Groups}},
  }

  if len(filter.Parties) > 0 {
    ledgerAdditionalFilter = append(ledgerAdditionalFilter, bson.M{"Name": bson.M{"$in": filter.Parties}})
  }

  // Search Filter
  var useLedgerSearch = filter.Filter.SearchKey == "party-wise" || filter.Filter.SearchKey == "group-wise"
  if len(filter.Filter.SearchText) > 0 && useLedgerSearch {
    searchKey := "Name"
    searchField := osMod.GetFieldBySearchKey(searchKey)
    searchFilter := utils.GenerateSearchFilter(filter.Filter.SearchText, searchField)
    ledgerAdditionalFilter = append(ledgerAdditionalFilter, bson.M{"$and": &searchFilter})
  }

  parties := ledgers.GetLedgers(companyId, ledgerFilter, ledgerAdditionalFilter)
  partyByName := utils.ToDict(parties, func(party ledgers.MetaLedger) string { return party.Name })

  istLocation, _ := time.LoadLocation("Asia/Kolkata")
  partySummary := make(map[string][]AgingOverview, 0)

  var mutex sync.Mutex

  batchFunc := func(partyChunk []ledgers.MetaLedger, wait *sync.WaitGroup) {
    defer wait.Done()
    partyNames := utils.Select(partyChunk, func(party ledgers.MetaLedger) string { return party.Name })
    billFilter := []bson.M{{"LedgerName": bson.M{"$in": partyNames}}}

    billDbFilter := filter
    billDbFilter.Filter.SortKey = "Name"
    bills := getBills(companyId, billDbFilter, &billFilter)

    if len(bills) < 1 {
      mutex.Lock()
      for _, partyName := range partyNames {
        partyInfo := partyByName[partyName]
        overview := AgingOverview{
          PartyName:     partyName,
          LedgerGroup:   partyInfo.Group,
          CreditDays:    partyInfo.CreditPeriod,
          CreditLimit:   partyInfo.CreditLimit,
          OpeningAmount: 0,
          ClosingAmount: 0,
        }
        partySummary[partyName] = []AgingOverview{overview}
      }
      mutex.Unlock()
    }

    today := time.Now().UTC()
    for _, bill := range bills {
      overview := AgingOverview{
        PartyName:     bill.LedgerName,
        BillNumber:    &bill.BillNumber,
        LedgerGroup:   *bill.LedgerGroupName,
        OpeningAmount: bill.OpeningBalance.Value,
        ClosingAmount: bill.ClosingBalance.Value,
        Above30:       0,
        Above60:       0,
        Above90:       0,
        Above120:      0,
      }

      billDate := bill.BillDate.In(istLocation)
      dueDate := billDate
      if bill.DueDate != nil {
        dueDate = (*bill.DueDate).In(istLocation)
      }

      diff := today.Sub(dueDate)
      days := int(diff.Hours() / 24)

      var delayDays uint16
      if days > 0 {
        delayDays = uint16(days)
      } else {
        delayDays = 0
      }
      overview.BillDate = &billDate
      overview.DueDate = &dueDate
      overview.DelayDays = &delayDays
      overview.IsAdvance = bill.IsAdvance

      // Update aging categories
      if useAgingRange {
        if 30 <= delayDays && delayDays < 60 {
          overview.Above30 = bill.ClosingBalance.Value
        } else if 60 <= delayDays && delayDays < 90 {
          overview.Above60 = bill.ClosingBalance.Value
        } else if 90 <= delayDays && delayDays < 120 {
          overview.Above90 = bill.ClosingBalance.Value
        } else if delayDays >= 120 {
          overview.Above120 = bill.ClosingBalance.Value
        }
      } else {
        if delayDays >= 30 {
          overview.Above30 = bill.ClosingBalance.Value
        }
        if delayDays >= 60 {
          overview.Above60 = bill.ClosingBalance.Value
        }
        if delayDays >= 90 {
          overview.Above90 = bill.ClosingBalance.Value
        }
        if delayDays >= 120 {
          overview.Above120 = bill.ClosingBalance.Value
        }
      }

      mutex.Lock()
      partySummary[bill.LedgerName] = append(partySummary[bill.LedgerName], overview)
      mutex.Unlock()
    }
  }

  totalParties := float64(len(parties))
  batchSize := math.Ceil(totalParties * 0.25)
  batchCount := math.Ceil(totalParties / batchSize)

  var wg sync.WaitGroup
  for batchNumber := 0; batchNumber < int(batchCount); batchNumber++ {
    wg.Add(1)
    startIndex := int(batchSize) * batchNumber
    endIndex := startIndex + int(batchSize)
    if endIndex > len(parties) {
      endIndex = len(parties)
    }
    partyChunk := parties[startIndex:endIndex]
    go batchFunc(partyChunk, &wg)
  }
  wg.Wait()

  var outstandingOverview []AgingOverview
  for _, partyLedger := range parties {
    overview := AgingOverview{
      PartyName:    partyLedger.Name,
      LedgerGroup:  partyLedger.Group,
      CreditLimit:  partyLedger.CreditLimit,
      CreditDays:   partyLedger.CreditPeriod,
      OpeningAmount: 0,
      ClosingAmount: 0,
      Above30: 0,
      Above60: 0,
      Above90: 0,
      Above120: 0,
      TotalBills:   0,
    }

    if summary, exist := partySummary[partyLedger.Name]; exist {
      var billCount int
      for _, bill := range summary {
        overview.OpeningAmount += bill.OpeningAmount
        billCount++
        overview.ClosingAmount += bill.ClosingAmount
        overview.Above30 += bill.Above30
        overview.Above60 += bill.Above60
        overview.Above90 += bill.Above90
        overview.Above120 += bill.Above120
      }
      overview.TotalBills = billCount
      overview.Bills = &summary
    }
    outstandingOverview = append(outstandingOverview, overview)
  }

  // Sorting by aging range
  sortAsc := utils.GetValueBySortOrder(filter.Filter.SortOrder) == 1
  switch filter.Filter.SortKey {
  case Above120Wise:
    utils.SortByField(outstandingOverview, "Above120", sortAsc)
  case Above90Wise:
    utils.SortByField(outstandingOverview, "Above90", sortAsc)
  case Above60Wise:
    utils.SortByField(outstandingOverview, "Above60", sortAsc)
  case Above30Wise:
    utils.SortByField(outstandingOverview, "Above30", sortAsc)
  default:
    utils.SortByField(outstandingOverview, "PartyName", sortAsc)
  }

  // Paginate result
  return utils.Paginate(outstandingOverview, int(filter.Filter.Batch.Limit), int(filter.Filter.Batch.Offset))
}

