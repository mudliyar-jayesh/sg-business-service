package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

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

type OsShareSettings struct {
    ID primitive.ObjectID `bson:"_id,omitempty"`
    CompanyId string `bson:"CompanyId"`
    CutOffDate string `bson:"CutOffDate"`
    //ShowItemDetails bool `bson:"ShowItemDetails"`
    //MinOsAmount float64 `bson:"MinOsAmount"`
    DueDays int `bson:"DueDays"`
    OverDueDays int `bson:"OverDueDays"`
    SendAllDue bool `bson:"SendAllDue"`
    SendDueOnly bool `bson:"SendDueOnly"`
    EmailSetting EmailSettings `bson:"EmailSetting"`
    AutoReminderInterval ReminderInterval `bson:"AutoReminderInterval"`
    ReminderIntervalDays int `bson:"ReminderIntervalDays"`
    //ReminderDaysBeforeDue int `bson:"ReminderDaysBeforeDue"`
}

