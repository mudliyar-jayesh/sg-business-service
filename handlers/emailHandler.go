package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"sg-business-service/models"
	"strings"
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
		message += k + ": " + v + "\r\n"
	}
	message += "\r\n" + config.Body

	allRecipients := append(config.To, append(config.Cc, config.Bcc...)...)

	fmt.Println("Before Send")
	err := smtp.SendMail(config.SmtpServer+":"+config.SmtpPort, auth, from, allRecipients, []byte(message))
	fmt.Println("After Send")
	if err != nil {
		fmt.Println("Error in send %v", err)
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
func WriteToTemplateFromString(templateContent string, data interface{}) string {
	// Parse the HTML template string
	tmpl, err := template.New("webpage").Parse(templateContent)
	if err != nil {
		fmt.Println("Could not parse template:", err)
		return ""
	}

	// Create a buffer to write the output to
	var buf bytes.Buffer

	// Execute the template and inject the data into the buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		fmt.Println("Could not write template:", err)
		return ""
	}

	// Return the HTML content as a string
	return buf.String()
}
