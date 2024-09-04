package handlers

import (
    "bytes"
    "html/template"
    "net/smtp"
    "strings"
    "sg-business-service/models"
    "fmt"
)

func SendEmail(config models.EmailSettings) error {

    from := "mudliyar.jayesh@gmail.com"
    //password := "grzj vqdz pceo xghm "
    password := "grzjvqdzpceoxghm "


    auth := smtp.PlainAuth("", from, password, config.SmtpServer)
    fmt.Println("Auth Complete")

    headers := make(map[string]string)
    headers["From"] = from
    headers["To"] = strings.Join(config.To, ", ")
    if len(config.Cc) > 0 {
        headers["Cc"] = strings.Join(config.Cc, ", ")
    }
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = "text/html; charset=\"UTF-8\""

    headers["Subject"] = config.Subject

    message := ""

    for k, v := range headers {
        message += k + ": " + v+ "\r\n"
    }
    message += "\r\n" + config.Body

    allRecipients :=append(config.To, append(config.Cc, config.Bcc...)...)

    fmt.Println("Before Send")
    err := smtp.SendMail(config.SmtpServer + ":" + config.SmtpPort, auth, from, allRecipients, []byte(message))
    fmt.Println("After Send")
    if err != nil {
        fmt.Println("Error in send ")
        return err
    }
    return nil

}

func WriteToTemplate(templatePath string, data interface{}) string {

    pageFormat, err := template.ParseFiles(templatePath)
    if err != nil {
        fmt.Println("Could not parse template")
    }

    var buf bytes.Buffer
    err = pageFormat.Execute(&buf, data)
    if err != nil {
        fmt.Println("Could not write template")
    }

    return buf.String()

}
