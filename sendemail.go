package sendemail

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v5"
)

const (
	MAILGUN_API_FILE = "/usr/local/etc/mailgun.txt"
	FROM_EMAIL       = "api@fake.com"
	TO_EMAIL         = "me@fake.com"
	DOMAIN           = "fake.com"
)

func Send(content string, subject string) error {
	// Get values from vars, else use constants
	sender := FROM_EMAIL // Default to constant
	recipient := TO_EMAIL
	if envSender := os.Getenv("MG_EMAIL_FROM"); envSender != "" {
		sender = envSender // Override with environment variable if set
	}
	if envRecipient := os.Getenv("MG_EMAIL_TO"); envRecipient != "" {
		recipient = envRecipient // Override with environment variable if set
	}

	domain := DOMAIN
	if envDomain := os.Getenv("MG_DOMAIN"); envDomain != "" {
		domain = envDomain // Override with environment variable if set
	}

	// Get Mailgun API Key
	var apiString string
	if _, err := os.Stat(MAILGUN_API_FILE); err == nil {
		// File exists, read from file
		privateAPIKeyBytes, err := os.ReadFile(MAILGUN_API_FILE)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return err
		}
		apiString = strings.TrimSpace(string(privateAPIKeyBytes))
	} else {
		// File does not exist, try environment variable
		apiString = os.Getenv("MAILGUN_API")
	}

	mg := mailgun.NewMailgun(apiString)

	// 2. Set up mailgun credentials and email details
	hostname, _err := os.Hostname()
	if _err != nil {
		fmt.Println(_err)
		os.Exit(1)
	}
	emailSubject := subject + " on " + hostname
	plainTextContent := string(content) // Use file content directly as plain text

	// The message object allows you to add attachments and Bcc recipients
	message := mailgun.NewMessage(domain, sender, emailSubject, plainTextContent, recipient)

	// 3. Send the email

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10-second timeout
	resp, email_err := mg.Send(ctx, message)
	if email_err != nil {
		log.Printf("The message content of the failed message: %v", message)
		log.Fatal(email_err)
		return email_err
	}

	fmt.Printf("Email sent successfully, resp: %s\n", resp.Message)
	return nil
}
