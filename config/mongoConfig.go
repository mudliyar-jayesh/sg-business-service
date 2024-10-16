package config

import (
	"os"
)

const (
	AppDb     string = "BMRM"
	TallyDb   string = "NewTallyDesktopSync"
	SummaryDb string = "Summary"
)

const (
	OsTemplate            string = "OsTemplates"
	ContactPerson         string = "ContactPerson"
	FollowUp              string = "FollowUp"
	Bill                  string = "Bills"
	DSP                   string = "DSP"
	Ledger                string = "Ledgers"
	LedgerGroup           string = "Groups"
	Voucher               string = "Vouchers"
	Item                  string = "StockItems"
	ItemGroup             string = "StockGroups"
	SyncInfo              string = "LastSyncInfo"
	OsSummary             string = "OutstandingSummary"
	CollectionActionables string = "CollectionActionables"
)

type MongoConfig struct {
	Uri string
}

func LoadMongoConfig() *MongoConfig {
	return &MongoConfig{
		Uri: getEnv("SG_MONGO", "mongodb://softgen:QWAmTnsdBUaTL2z@118.139.167.125:27017/"),
	}
}

func LoadEmailTemplate() string {
	return getEnv("SG_Template", "")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
