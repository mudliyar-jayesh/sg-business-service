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
)

func (r ReminderInterval) String() string {
    return [...]string{"Daily", "Weekly", "Monthly", "DayWise"}[r-1]
}

type OsShareSettings struct {
    ID primitive.ObjectID `bson:"_id,omitempty"`
    CompanyId string `bson:"CompanyId"`
    CutOffDate string `bson:"CutOffDate"`
    SendAllDue bool `bson:"SendAllDue"`
    SendDueOnly bool `bson:"SendDueOnly"`
    EmailSetting EmailSettings `bson:"EmailSetting"`
    AutoReminderInterval ReminderInterval `bson:"AutoReminderInterval"`
    ReminderIntervalDays int `bson:"ReminderIntervalDays"`
}

