package internal

import (
	"imap-sync/config"
	"net/smtp"
	"strings"
)

var (
	smtpHost     string
	smtpPort     string
	from         string
	smtpUsername string
	smtpPassword string
)

func Notify(task *Task, isSuccessful bool) {
	smtpHost = config.Conf.Email.SMTPHost
	smtpPort = config.Conf.Email.SMTPPort
	from = config.Conf.Email.SMTPFrom
	smtpUsername = config.Conf.Email.SMTPUser
	smtpPassword = config.Conf.Email.SMTPPassword

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

	auth := smtp.CRAMMD5Auth(smtpUsername, smtpPassword)

	msg := []byte("From: " + from + "\r\n" +
		"To: " + toHeader + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		message + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, accounts, msg)
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}
