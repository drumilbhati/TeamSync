package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
)

func SendOTP(userEmail, userName, otp string) error {
	from := os.Getenv("FROM_MAIL")
	password := os.Getenv("PASS_MAIL")

	if from == "" || password == "" || from == "your-email@gmail.com" {
		return fmt.Errorf("email credentials are not set in .env file. Please configure FROM_MAIL and PASS_MAIL")
	}

	to := []string{userEmail}

	host := "smtp.gmail.com"

	port := "587"

	subject := "Subject: Your TeamSync Verification Code\n"

	mime := "MIME-version:1.0;\nContent-Type: text/plain;charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf("Hi %s, \n\nYour verification code is: %s,\n\nThis is valid for 10 minutes.", userName, otp)

	msg := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", from, password, host)

	fmt.Printf("Attempting to send email from %s to %s via %s:%s...\n", from, userEmail, host, port)
	err := smtp.SendMail(host+":"+port, auth, from, to, msg)
	if err != nil {
		fmt.Printf("SMTP error occurred: %v\n", err)
		return fmt.Errorf("SMTP error: %w", err)
	}

	fmt.Printf("SMTP call completed successfully for %s\n", userEmail)
	return nil
}

func GenerateRandomNumber() string {
	// Generate a random number between 100000 and 999999 (inclusive)
	num := rand.Intn(900000) + 100000
	return strconv.Itoa(num)
}
