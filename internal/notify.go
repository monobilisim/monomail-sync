package internal

import (
	"flag"
	"net/smtp"
	"strings"
)

var (
	smtpHost     = flag.String("smtpHost", "", "SMTP server")
	smtpPort     = flag.String("smtpPort", "", "SMTP port")
	from         = flag.String("from", "", "SMTP mail sender")
	smtpUsername = flag.String("smtpUser", "", "SMTP username")
	smtpPassword = flag.String("smtpPass", "", "SMTP password")
)

func Notify(task *Task, isSuccessful bool) {
	if !isSuccessful {
		err := sendMmail([]string{task.SourceAccount, task.DestinationAccount}, Data["notify"]["fail"], Data["notify"]["fail_msg"])
		if err != nil {
			log.Info(err)
		}
	} else {
		err := sendMmail([]string{task.SourceAccount, task.DestinationAccount}, Data["notify"]["success"], Data["notify"]["success_msg"])
		if err != nil {
			log.Info(err)
		}
	}
}

func sendMmail(accounts []string, status string, text string) error {
	toHeader := strings.Join(accounts, ",")
	subject := "Monomail-sync " + status
	message := accounts[0] + " - " + accounts[1] + text

	auth := smtp.CRAMMD5Auth(*smtpUsername, *smtpPassword)

	msg := []byte("From: " + *from + "\r\n" +
		"To: " + toHeader + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		message + "\r\n")

	err := smtp.SendMail(*smtpHost+":"+*smtpPort, auth, *from, accounts, msg)
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}
