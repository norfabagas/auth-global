package smtp

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

var config_email = os.Getenv("CONFIG_SMTP_EMAIL")
var config_password = os.Getenv("CONFIG_SMTP_PASSWORD")
var config_host = os.Getenv("CONFIG_SMTP_HOST")
var config_port = os.Getenv("CONFIG_SMTP_PORT")

func Send(to, cc []string, subject, message string) error {
	body := "From: " + config_email + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth(
		"",
		config_email,
		config_password,
		config_port,
	)
	smtpAddr := fmt.Sprintf("%s:%s",
		config_host,
		config_port,
	)

	err := smtp.SendMail(
		smtpAddr,
		auth,
		config_email,
		append(to, cc...),
		[]byte(body),
	)

	if err != nil {
		return err
	}

	return nil
}
