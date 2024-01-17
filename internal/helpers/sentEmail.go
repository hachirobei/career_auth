package helpers

import (
    "net/smtp"
    "fmt"
    "log"
    "os"
)

// SendPasswordResetEmail sends a password reset email to the specified recipient
func SendPasswordResetEmail(recipientEmail, token string) error {
    smtpServer := os.Getenv("SERVERSMTP")
    port := os.Getenv("PORTSMTP")
    senderEmail := os.Getenv("EMAILSMTP")
    password := os.Getenv("PASSWORDSMTP") 

    // Set up authentication information.
    auth := smtp.PlainAuth("", senderEmail, password, smtpServer)

    // Construct the email message
    subject := "Subject: Password Reset Instructions\n"
    mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
    body := fmt.Sprintf("<html><body>You requested a password reset. Please use the following token: <strong>%s</strong></body></html>", token)
    msg := []byte(subject + mime + body)

    // Sending email
    err := smtp.SendMail(smtpServer+":"+port, auth, senderEmail, []string{recipientEmail}, msg)
    if err != nil {
        log.Printf("SMTP Error: %s\n", err)
        return err
    }

    log.Println("Email sent successfully")
    return nil
}
