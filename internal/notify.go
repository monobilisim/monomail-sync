package internal

import (
	"net/smtp"
)

func Email(creds Credentials) error {

	smtpHost := ""
	smtpPort := ""
	from := ""
	username := ""
	password := ""
	to := creds.Account
	subject := ""
	message := ""

	auth := smtp.CRAMMD5Auth(username, password)

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: [" + "" + "] " + subject + "\r\n\r\n" +
		message + "\r\n")

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}
