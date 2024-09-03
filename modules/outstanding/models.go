package outstanding

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
    SearchKey string
    SortKey string
    SortOrder string
    ReportOnType ReportType
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
