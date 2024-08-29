package handlers

import (
    "net/smtp"
    "strings"
    "sg-business-service/models"
)

func SendEmail(config models.EmailSettings) error {

    from := "mudliyar.jayesh@gmail.com"
    password := ""

    auth := smtp.PlainAuth("", from, password, config.SmtpServer)

    headers := make(map[string]string)
    headers["From"] = from
    headers["To"] = strings.Join(config.To, ", ")
    if len(config.Cc) > 0 {
        headers["Cc"] = strings.Join(config.Cc, ", ")
    }

    headers["Subject"] = config.Subject

    message := ""

    for k, v := range headers {
        message += k + ": " + v+ "\r\n"
    }
    message += "\r\n" + config.Body

    allRecipients :=append(config.To, append(config.Cc, config.Bcc...)...)

    err := smtp.SendMail(config.SmtpServer + ":" + config.SmtpPort, auth, from, allRecipients, []byte(message))
    if err != nil {
        return err
    }
    return nil

}
