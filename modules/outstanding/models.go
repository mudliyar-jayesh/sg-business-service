package outstanding

import (
	"sg-business-service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

type LocationReportType string

const (
	StateWise    LocationReportType = "State"
	RegionWise   LocationReportType = "Region"
	DistrictWise LocationReportType = "District"
	PincodeWise  LocationReportType = "Pincode"
)

type OsReportFilter struct {
	PartyName    string
	SearchText   string
	Limit        int64
	Offset       int64
	Groups       []string
	DueFilter    DueDayFilter
	SearchKey    string
	SortKey      string
	SortOrder    string
	ReportOnType ReportType
}

type OsLocationFilter struct {
	Filter       models.RequestFilter
	State        string
	Locations    []string
	LocationType LocationReportType
}

type LocationOverview struct {
	LocationName  string
	OpeningAmount float64
	ClosingAmount float64
}

type MetaBill struct {
	BillNumber    string                  `bson:"Name"`
	PartyName     string                  `bson:"LedgerName"`
	ParentGroup   *string                 `bson:"LedgerGroupName"`
	PendingAmount *models.FloatFromString `bson:"Amount"`
	OpeningAmount *models.FloatFromString `bson:"OpeningAmount"`
	BillDate      *primitive.DateTime     `bson:"BillDate"`
	DueDate       *primitive.DateTime     `bson:"DueDate"`
}

type Bill struct {
	LedgerName      string
	LedgerGroupName string
	BillName        string
	BillDate        string
	DueDate         string
	DelayDays       int32

	OpeningAmount float64
	ClosingAmount float64
	Amount        float64
	DueAmount     float64
	OverDueAmount float64

	PendingPercentage float64
	PaidPercentage    float64

	AmountStr        string
	OpeningAmountStr string
	ClosingAmountStr string
	OverDueAmountStr string
}

type OsShareSettings struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CompanyId    string             `bson:"CompanyId"`
	CutOffDate   string             `bson:"CutOffDate"`
	TemplateName *string            `bson:"TemplateName"`
	//ShowItemDetails bool `bson:"ShowItemDetails"`
	//MinOsAmount float64 `bson:"MinOsAmount"`
	DueDays              int                  `bson:"DueDays"`
	OverDueDays          int                  `bson:"OverDueDays"`
	SendAllDue           bool                 `bson:"SendAllDue"`
	SendDueOnly          bool                 `bson:"SendDueOnly"`
	SendOverDueOnly      bool                 `bson:"SendOverDueOnly"`
	EmailSetting         models.EmailSettings `bson:"EmailSetting"`
	AutoReminderInterval ReminderInterval     `bson:"AutoReminderInterval"`
	ReminderIntervalDays int                  `bson:"ReminderIntervalDays"`
	//ReminderDaysBeforeDue int `bson:"ReminderDaysBeforeDue"`
}

type ReminderInterval int

const (
	Daily ReminderInterval = iota
	Weekly
	Monthly
	DayWise
	// FortNightly
)

func (r ReminderInterval) String() string {
	return [...]string{"Daily", "Weekly", "Monthly", "DayWise"}[r-1]
}

type PartyOverview struct {
	PartyName    string
	BillNumber   string
	TotalBills   int32
	TotalOpening float64
	TotalClosing float64
}
