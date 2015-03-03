package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
)

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func main() {
	// Set up authentication information.
//	smtpServer := "192.241.204.134"
    smtpServer := "127.0.0.1"
	auth := smtp.PlainAuth(
		"",
		"",
		"",
		smtpServer,
	)

	from := mail.Address{"", "+14257772300@smtp.zipwhip.com"}
	to := mail.Address{"Cell", "+12068597896@smtp.zipwhip.com"}

	//title := "8ef1211f-d9f2-4c81-906f-7d27da5a32f8:309626613"
    title := "+14257772300@smtp.ziphip.com"
	body := "Hello World\n\nHow are you today?\n\nGood Thanks!"

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["Zipwhip-Auth"] = "asdfasdf"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		//smtpServer+":25",
        smtpServer+":10025",
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	//[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
}
