package models

type EmailBodyType int

const (
    TextBody EmailBodyType = iota
    HtmlBody
)


func (r EmailBodyType) String() string {
    return [...]string{"TextBody", "HtmlBody"}[r-1]
}

type EmailSettings struct {
    To string `bson:"To"`
    Cc string `bson:"Cc"`
    Bcc string `bson:"Bcc"`
    Subject string `bson:"Subject"`
    Body string `bson:"Body"`
    BodyType EmailBodyType `bson:"BodyType"`
    Signature string `bson:"Signature"`
}



