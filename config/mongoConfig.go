package config

import (
	"os"
)

const (
    AppDb string = "BMRM"
    TallyDb string = "NewTallyDesktopSync"
)

const (
    Bill string = "Bills"
    Ledger string = "Ledgers"
    LedgerGroup string = "Groups"
    Voucher string = "Vouchers"
    Item string = "StockItems"
    ItemGroup string = "StockGroups"
    ContactPerson string = "ContactPerson"
    FollowUp string = "FollowUp"
)


type MongoConfig struct {
    Uri string
}

func LoadMongoConfig() *MongoConfig {
    return &MongoConfig {
        Uri: getEnv("SG_MONGO", "mongodb://softgen:QWAmTnsdBUaTL2z@118.139.167.125:27017/"),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
