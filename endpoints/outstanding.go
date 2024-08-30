package endpoints
import (
    "time"
    "context"
    "net/http"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "sg-business-service/handlers"
    "sg-business-service/utils"
    "strconv"
)


type DueDayFilter int
type ReportType int

const (
    PartyWise ReportType = iota
    BillWise
)
const (
    AllBills DueDayFilter = iota
    PendingBills
    DueBills
    OverDueBills
)

type OsReportFilter struct {
    PartyName string
    SearchText string
    Limit int64
    Offset int64
    Groups []string
    DueFilter DueDayFilter
    DueDays int
    OverDueDays int
    SearchKey string
    SortKey string
    SortOrder string
    ReportOnType ReportType
}

func getFieldBySortKey(sortKey string) string {
    var fieldBySortKey = make(map[string]string)
    fieldBySortKey["Party"] = "LedgerName"
    fieldBySortKey["Group"] = "LedgerGroupName"
    fieldBySortKey["Bill"] = "Name"

    sortField, exists := fieldBySortKey[sortKey]
    if exists {
        return sortField
    }
    return fieldBySortKey["Party"]
}
func getFieldBySearchKey(searchKey string) string {
    var fieldBySearchKey = make(map[string]string)
    fieldBySearchKey["Party"] = "LedgerName"
    fieldBySearchKey["Group"] = "LedgerGroupName"
    fieldBySearchKey["Bill"] = "Name"

    searchField, exists := fieldBySearchKey[searchKey]
    if exists {
        return searchField
    }
    return fieldBySearchKey["Party"]
}

func SearchLedgers(res http.ResponseWriter, req *http.Request) {
    companyId := req.Header.Get("CompanyId")
    searchKey := req.URL.Query().Get("searchKey")

    var collection = handlers.GetCollection("NewTallyDesktopSync", "Ledgers")
    var mongoHandler = handlers.NewMongoHandler(collection)

    var filter = bson.M {
        "GUID": bson.M {
            "$regex": "^"+companyId,
        },
    }

    var page int64= 25
    var usePagination = true
    if len(searchKey) > 0 {
        page = 0
        usePagination = false
        filter["$and"] = utils.GenerateSearchFilter(searchKey, "Name")
    }

    docFilter := handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter:filter,
        UsePagination: usePagination,
        Limit: page,
        Offset: 0,
        Projection: bson.M{
            "Name": 1,
            "_id": 0,
        },
    }

    var results handlers.DocumentResponse= mongoHandler.FindDocuments(docFilter)

    response := utils.NewResponseStruct(results.Data, len(results.Data))
    response.ToJson(res)
}

func GetCachedGroups(res http.ResponseWriter, req *http.Request) {
    companyId := req.Header.Get("CompanyId")
    isDebit := utils.GetBoolFromQuery(req, "isDebit")

    var parentName string = "Current Assets"
    if !isDebit {
        parentName = "Current Liabilities"
    }

    var groups = handlers.CachedGroups.GetChildrenNames(companyId, parentName)

    response := utils.NewResponseStruct(groups, len(groups))
    response.ToJson(res)
}


func GetOutstandingReport(res http.ResponseWriter, req *http.Request) {
    var collection = handlers.GetCollection("NewTallyDesktopSync", "Bills")
    var mongoHandler = handlers.NewMongoHandler(collection)


    companyId := req.Header.Get("CompanyId")
    isDebit := utils.GetBoolFromQuery(req, "isDebit")

    var parentName string = "Current Assets"
    if !isDebit {
        parentName = "Current Liabilities"
    }

    var groups = handlers.CachedGroups.GetChildrenNames(companyId, parentName)

    reqBody, err := utils.ReadRequestBody[OsReportFilter](req)
    if err != nil {
        http.Error(res, "Unable to read request body", http.StatusBadRequest)
        return
    }

    if len(reqBody.Groups) > 0 {
        groups = utils.Intersection(groups, reqBody.Groups)
    }

    var filter = bson.M {
        "CompanyId": companyId,
        "LedgerGroupName": bson.M {
            "$in": groups,
        },
        "ClosingBal": bson.M {"$ne": nil },
        "ClosingBal.IsDebit": isDebit,
    }

    if len(reqBody.PartyName) > 0 {
        filter["LedgerName"] = reqBody.PartyName
    } else if len(reqBody.SearchText) > 0 {
        var searchField = getFieldBySearchKey(reqBody.SearchKey)
        filter["$and"] = utils.GenerateSearchFilter(reqBody.SearchText, searchField)
    }

    var usePagination = reqBody.Limit != 0
    if reqBody.ReportOnType == PartyWise {
        usePagination = false
    }

    docFilter := handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter:filter,
        UsePagination: usePagination,
        Limit: reqBody.Limit,
        Offset: reqBody.Offset,
        Projection: bson.M{
            "LedgerName": 1,
            "LedgerGroupName": 1,
            "BillDate": "$BillDate.Date",
            "DueDate": "$BillCreditPeriod.DueDate",
            "Amount": "$ClosingBal.Amount",
            "Name": "$Name",
            "_id": 0,
        },
        Sorting: bson.D {
            {
                Key: getFieldBySortKey(reqBody.SortKey),
                Value: utils.GetValueBySortOrder(reqBody.SortOrder),
            },
        },
    }

    var results handlers.DocumentResponse= mongoHandler.FindDocuments(docFilter)

    settingCollection := handlers.GetCollection("BMRM", "OutstandingSettings")
    var settingsHandler = handlers.NewMongoHandler(settingCollection)

    docFilter = handlers.DocumentFilter {
        Ctx: context.TODO(),
        Filter: bson.M {
            "CompanyId": companyId,
        },
        UsePagination: false,
        Limit: 0,
        Offset: 0,
    }

    settings := settingsHandler.FindDocuments(docFilter)
    if settings.Err != nil {
        http.Error(res, "No Data", http.StatusBadRequest)
        return;
    }

    var dueDays int32
    var overDueDays int32

    if len(settings.Data) > 0 {
        dueDays = settings.Data[0]["DueDays"].(int32)
        overDueDays= settings.Data[0]["OverDueDays"].(int32)
    }
    var bills []Bill
    istLocation, _ := time.LoadLocation("Asia/Kolkata")
    for _, item := range results.Data {
        billDateValue := item["BillDate"].(primitive.DateTime).Time()
        billDate := billDateValue.In(istLocation).Format("2006-01-02 15:04:05")

        var dueDate string
        dueDateValue := item["DueDate"]
        if dueDateValue == nil {
            dueDate = billDate
        } else {
           dueDate = dueDateValue.(primitive.DateTime).Time().In(istLocation).Format("2006-01-02 15:04:05")
        }

        layout := "2006-01-02 15:04:05"
        parsedTime, _:= time.Parse(layout, dueDate)
        today := time.Now().UTC()

        // Calculate the difference
        diff := today.Sub(parsedTime)

        // Get the difference in days
        days := int32(diff.Hours() / 24)

        var bill Bill = Bill {
            LedgerName: item["LedgerName"].(string),
            LedgerGroupName: item["LedgerGroupName"].(string),
            BillName: item["Name"].(string),
            DueDate: dueDate,
            BillDate: billDate,
            DelayDays: days,
        }
        var amount = parseFloat64(item["Amount"])
        var dueFilter = AllBills
        if days >= dueDays && days <= overDueDays {
            bill.DueAmount = amount
            dueFilter = DueBills
        } else if days > overDueDays {
            bill.OverDueAmount = amount
            dueFilter = OverDueBills
        } else {
            bill.Amount = amount
            dueFilter = PendingBills
        }

        if reqBody.DueFilter != AllBills && reqBody.DueFilter != dueFilter{
            continue
        }
        bills = append(bills, bill)
    }

    if reqBody.ReportOnType == PartyWise {
        var groupedBills = utils.GroupByKey(bills, "LedgerName")

        var partyBills []Bill
        for _, group := range groupedBills {
            if len(group) < 1 {
                continue;
            }

            var firstEntry = group[0]

            var totalAmount float64 = 0
            var totalDue float64 = 0
            var totalOverDue float64 = 0

            for _, bill := range group {
                totalAmount += bill.Amount
                totalDue += bill.DueAmount
                totalOverDue += bill.OverDueAmount
            }

            var partyBill Bill = Bill {
                LedgerName: firstEntry.LedgerName,
                LedgerGroupName: firstEntry.LedgerGroupName,
                BillName: "",
                DueDate: "",
                BillDate: "",
                DelayDays: 0,
                Amount:totalAmount,
                DueAmount: totalDue,
                OverDueAmount: totalOverDue,
            }

            partyBills = append(partyBills, partyBill)
        }


        // Skip the first 8 records
        skip := reqBody.Offset * reqBody.Limit
        length := int64(len(partyBills))
        if skip > length {
            skip = 0
        }
        skipped := partyBills[skip:]

        // Take the next 5 records
        take := skip + reqBody.Limit
        if take > length {
            take = length
        }
        takenBills := skipped[:take]

        response := utils.NewResponseStruct(takenBills, len(takenBills))
        response.ToJson(res)
        return 
    }

    response := utils.NewResponseStruct(bills, len(bills))
    response.ToJson(res)
}


func TempOS(res http.ResponseWriter, req *http.Request) {
    var collection = handlers.GetCollection("NewTallyDesktopSync", "Bills")
    var mongoHandler = handlers.NewMongoHandler(collection)
    companyId := req.Header.Get("CompanyId")

    pipeline := mongo.Pipeline{
        {  {"$match", bson.D{
            {"VoucherId", bson.M {
                "$regex": "^"+ companyId,
            }},
        }}},
        {{"$lookup", bson.D{
            {"from", "Vouchers"},
            {"localField", "VoucherId"},
            {"foreignField", "GUID"},
            {"as", "voucher_info"},
        }}},
        {{"$unwind", "$voucher_info"}},
        {
            {
                "$project", bson.D {
                    {"BillName", 1 },
                    {"SalesPerson", 1 },
                    {"BrokerName", 1 },
                    {"TotalBrokerCommission", 1 },
                    {"TotalSalesIncentive", 1 },
                    {"Ledger", bson.D{
                        {"Name", bson.D{
                            {"$arrayElemAt", bson.A{"$voucher_info.Ledgers.LedgerName", 0}},
                        }},
                        {"Amount", bson.D{
                            {"$arrayElemAt", bson.A{"$voucher_info.Ledgers.Amount.Amount", 0}},
                        }},
                    }},
                }}},
            }

    results := mongoHandler.AggregatePipeline("NewTallyDesktopSync", "VoucherUserTypeDetails", pipeline)
    response := utils.NewResponseStruct(results, len(results))
    response.ToJson(res)
}

func parseFloat64(value interface{}) float64 {
    var result float64
    switch v := value.(type) {
    case float64:
        result = v
    case int:
        result = float64(v)
    case string:
        parsed, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return 0 // Return default value on error
        }
        result = parsed
    default:
        return 0 // Return default value if type is not handled
    }
    return result
}

type Temp struct {
    Data interface{}
    Count int
}
type Bill struct {
    LedgerName string 
    LedgerGroupName string 
    BillName string
    BillDate string
    DueDate string
    DelayDays int32
    Amount float64
    DueAmount float64
    OverDueAmount float64
}
