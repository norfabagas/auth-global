package smtp

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func Send(to, cc []string, subject, message string) error {
	body := "From: " + os.Getenv("CONFIG_SMTP_EMAIL") + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth(
		"",
		os.Getenv("CONFIG_SMTP_EMAIL"),
		os.Getenv("CONFIG_SMTP_PASSWORD"),
		os.Getenv("CONFIG_SMTP_PORT"),
	)
	smtpAddr := fmt.Sprintf("%s:%s",
		os.Getenv("CONFIG_SMTP_HOST"),
		os.Getenv("CONFIG_SMTP_PORT"),
	)

	err := smtp.SendMail(
		smtpAddr,
		auth,
		os.Getenv("CONFIG_SMTP_EMAIL"),
		append(to, cc...),
		[]byte(body),
	)

	if err != nil {
		return err
	}

	return nil
}
