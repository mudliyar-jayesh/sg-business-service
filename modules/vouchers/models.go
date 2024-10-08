
package vouchers

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetaVoucher struct {
	GUID                  string             `bson:"GUID"`
	Date                  DateField          `bson:"Date"`
	VoucherType           string             `bson:"VoucherType"`
	ReferenceDate         DateField          `bson:"ReferenceDate"`
	Reference             string             `bson:"Reference"`
  InventoryAllocations  []Inventory        `bson:"InventoryAllocations"`
	Ledgers               []Ledger           `bson:"Ledgers"`
}

type Voucher struct {
	ID                    primitive.ObjectID `bson:"_id"`
	Action                int                `bson:"Action"`
	MasterID              int64              `bson:"MasterId"`
	GUID                  string             `bson:"GUID"`
	RemoteID              string             `bson:"RemoteId"`
	AlterID               int64              `bson:"AlterId"`
	Date                  DateField          `bson:"Date"`
	ReferenceDate         DateField          `bson:"ReferenceDate"`
	Reference             string             `bson:"Reference"`
	VoucherType           string             `bson:"VoucherType"`
	VoucherTypeID         string             `bson:"VoucherTypeId"`
	View                  int                `bson:"View"`
	VoucherGSTClass       string             `bson:"VoucherGSTClass"`
	IsCostCentre          BoolField          `bson:"IsCostCentre"`
	CostCentreName        string             `bson:"CostCentreName"`
	VoucherEntryMode      string             `bson:"VoucherEntryMode"`
	IsInvoice             BoolField          `bson:"IsInvoice"`
	VoucherNumber         string             `bson:"VoucherNumber"`
	IsOptional            BoolField          `bson:"IsOptional"`
	EffectiveDate         DateField          `bson:"EffectiveDate"`
	Narration             string             `bson:"Narration"`
	PriceLevel            string             `bson:"PriceLevel"`
	BillToPlace           string             `bson:"BillToPlace"`
	IRN                   string             `bson:"IRN"`
	IRNAckNo              string             `bson:"IRNAckNo"`
	IRNAckDate            string             `bson:"IRNAckDate"`
	DeliveryNoteNo        string             `bson:"DeliveryNoteNo"`
	ShippingDate          string             `bson:"ShippingDate"`
	DispatchFromName      string             `bson:"DispatchFromName"`
	DispatchFromStateName string             `bson:"DispatchFromStateName"`
	DispatchFromPinCode   string             `bson:"DispatchFromPinCode"`
	DispatchFromPlace     string             `bson:"DispatchFromPlace"`
	DeliveryNotes         DeliveryNotes      `bson:"DeliveryNotes"`
	DispatchDocNo         string             `bson:"DispatchDocNo"`
	BasicShippedBy        string             `bson:"BasicShippedBy"`
	Destination           string             `bson:"Destination"`
	CarrierName           string             `bson:"CarrierName"`
	BillofLandingNo       string             `bson:"BillofLandingNo"`
	BillofLandingDate     string             `bson:"BillofLandingDate"`
	PlaceOfReceipt        string             `bson:"PlaceOfReceipt"`
	ShipOrFlightNo        string             `bson:"ShipOrFlightNo"`
	LandingPort           string             `bson:"LandingPort"`
	DischargePort         string             `bson:"DischargePort"`
	DesktinationCountry   string             `bson:"DesktinationCountry"`
	ShippingBillNo        string             `bson:"ShippingBillNo"`
	ShippingBillDate      string             `bson:"ShippingBillDate"`
	PortCode              string             `bson:"PortCode"`
	BasicDueDateofPayment string             `bson:"BasicDueDateofPayment"`
	OrderReference        string             `bson:"OrderReference"`
	PartyName             string             `bson:"PartyName"`
	PartyLedgerID         string             `bson:"PartyLedgerId"`
	GSTRegistration       GSTRegistration    `bson:"GSTRegistration"`
	VoucherNumberSeries   string             `bson:"VoucherNumberSeries"`
	PartyMailingName      string             `bson:"PartyMailingName"`
	State                 string             `bson:"State"`
	Country               string             `bson:"Country"`
	RegistrationType      string             `bson:"RegistrationType"`
	PartyGSTIN            string             `bson:"PartyGSTIN"`
	PlaceOfSupply         string             `bson:"PlaceOfSupply"`
	PINCode               string             `bson:"PINCode"`
	ConsigneeName         string             `bson:"ConsigneeName"`
	ConsigneeMailingName  string             `bson:"ConsigneeMailingName"`
	ConsigneeState        string             `bson:"ConsigneeState"`
	ConsigneeCountry      string             `bson:"ConsigneeCountry"`
	ConsigneeGSTIN        string             `bson:"ConsigneeGSTIN"`
	ConsigneePinCode      string             `bson:"ConsigneePinCode"`
	Address               []string           `bson:"Address"`
	BuyerAddress          []string           `bson:"BuyerAddress"`
	IsCancelled           BoolField          `bson:"IsCancelled"`
	OverrideEWayBill      BoolField          `bson:"OverrideEWayBillApplicability"`
	EWayBillDetails       string             `bson:"EWayBillDetails"`
	Ledgers               []Ledger           `bson:"Ledgers"`
	InventoryAllocations  []Inventory        `bson:"InventoryAllocations"`
	InventoriesOut        []string           `bson:"InventoriesOut"`
	InventoriesIn         []string           `bson:"InventoriesIn"`
	CategoryEntry         string             `bson:"CategoryEntry"`
	AttendanceEntries     string             `bson:"AttendanceEntries"`
	Dt                    string             `bson:"Dt"`
	VchType               string             `bson:"VchType"`
	MasterId              string             `bson:"_MasterId"`
}

type DateField struct {
	Date time.Time `bson:"$date"`
}

type BoolField struct {
	Value bool `bson:"Value"`
}

type DeliveryNotes struct {
	ShippingDate string `bson:"ShippingDate"`
	DeliveryNote string `bson:"DeliveryNote"`
}

type GSTRegistration struct {
	TaxType          string `bson:"TaxType"`
	TaxRegistration  string `bson:"TaxRegistration"`
	RegistrationName string `bson:"RegistrationName"`
}

type Ledger struct {
	IndexNumber          int                `bson:"IndexNumber"`
	LedgerName           string             `bson:"LedgerName"`
	LedgerID             string             `bson:"LedgerId"`
	LedgerTaxType        string             `bson:"LedgerTaxType"`
	LedgerType           string             `bson:"LedgerType"`
	IsDeemedPositive     BoolField          `bson:"IsDeemedPositive"`
	Amount               AmountField        `bson:"Amount"`
	CostCategoryAllocations []string        `bson:"CostCategoryAllocations"`
	AdAllocType          int                `bson:"AdAllocType"`
	IsPartyLedger        BoolField          `bson:"IsPartyLedger"`
	SWIFTCode            string             `bson:"SWIFTCode"`
	BankAllocations      string             `bson:"BankAllocations"`
	BillAllocations      []BillAllocation   `bson:"BillAllocations"`
	InventoryAllocations []string           `bson:"InventoryAllocations"`
}

type AmountField struct {
	Amount       string `bson:"Amount"`
	ForexAmount  string `bson:"ForexAmount"`
	RateOfExchange string `bson:"RateOfExchange"`
	Currency     string `bson:"Currency"`
	IsDebit      bool   `bson:"IsDebit"`
	PreserveAmount bool `bson:"PreserveAmount"`
}

type BillAllocation struct {
	BillType           int        `bson:"BillType"`
	Name               string     `bson:"Name"`
	BillID             int        `bson:"BillId"`
	BillCreditPeriod   BillPeriod `bson:"BillCreditPeriod"`
	Amount             AmountField `bson:"Amount"`
}

type BillPeriod struct {
	BillDate  DateField `bson:"BillDate"`
	Value     int       `bson:"Value"`
	Suffix    int       `bson:"Suffix"`
	DueDate   DateField `bson:"DueDate"`
}

type Inventory struct {
	UserDescriptions    []string       `bson:"UserDescriptions"`
	IndexNumber         int            `bson:"IndexNumber"`
	StockItemName       string         `bson:"StockItemName"`
	StockItemID         string         `bson:"StockItemId"`
	BOMName             string         `bson:"BOMName"`
	IsScrap             BoolField      `bson:"IsScrap"`
	IsDeemedPositive    BoolField      `bson:"IsDeemedPositive"`
	Rate                RateField      `bson:"Rate"`
	ActualQuantity      QuantityField  `bson:"ActualQuantity"`
	BilledQuantity      QuantityField  `bson:"BilledQuantity"`
	Amount              AmountField    `bson:"Amount"`
	BatchAllocations    []BatchAllocation `bson:"BatchAllocations"`
	CostCategoryAllocations []string   `bson:"CostCategoryAllocations"`
	Ledgers             []Ledger       `bson:"Ledgers"`
}

type RateField struct {
	RatePerUnit string `bson:"RatePerUnit"`
	Unit        string `bson:"Unit"`
}

type QuantityField struct {
	Number       string    `bson:"Number"`
	PrimaryUnits UnitField `bson:"PrimaryUnits"`
	SecondaryUnits string  `bson:"SecondaryUnits"`
}

type UnitField struct {
	Number string `bson:"Number"`
	Unit   string `bson:"Unit"`
}

type BatchAllocation struct {
	ManufacturedOn string  `bson:"ManufacturedOn"`
	TrackingNo     string  `bson:"TrackingNo"`
	OrderNo        string  `bson:"OrderNo"`
	GodownName     string  `bson:"GodownName"`
	BatchName      string  `bson:"BatchName"`
	OrderDueDate   BillPeriod `bson:"OrderDueDate"`
	Amount         AmountField `bson:"Amount"`
	ActualQuantity QuantityField `bson:"ActualQuantity"`
	BilledQuantity QuantityField `bson:"BilledQuantity"`
}

