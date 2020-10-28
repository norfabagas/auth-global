package smtp

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

// Run this function asynchronously to maintain performance
func Send(to []string, cc []string, subject, message string) {
	email := os.Getenv("CONFIG_SMTP_EMAIL")
	password := os.Getenv("CONFIG_SMTP_PASSWORD")
	host := os.Getenv("CONFIG_SMTP_HOST")
	port := os.Getenv("CONFIG_SMTP_PORT")

	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.Println(err)
	}

	body := "From: " + email + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", email, password, host)
	smtpAddr := fmt.Sprintf("%s:%d", host, intPort)

	err = smtp.SendMail(smtpAddr, auth, email, append(to, cc...), []byte(body))
	if err != nil {
		log.Println(err)
	}
}
