package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendOTP(userEmail, userName, otp string) error {
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
	if fromEmail == "" {
		return fmt.Errorf("SENDGRID_FROM_EMAIL not set in .env")
	}

	from := mail.NewEmail("TeamSync", fromEmail)

	to := mail.NewEmail(userName, userEmail)

	plainTextContent := fmt.Sprintf("Hi %s, your verification code for TeamSync is: %v\n This code will expire in 10 minutes", userName, otp)
	htmlContent := fmt.Sprintf("<strong>Hi %s,</strong><br><p>Your verification code for TeamSync is:</p><h1>%s</h1><p>This code will expire in 10 minutes.</p>", userName, otp)

	message := mail.NewSingleEmail(from, "Your TeamSync Verification Code", to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode != 202 {
		return fmt.Errorf("sendgrid api returned non-success status: %d. Body: %s", response.StatusCode, response.Body)
	}

	return nil
}

func GenerateRandomNumber() string {
	// Generate a random number between 100000 and 999999 (inclusive)
	num := rand.Intn(900000) + 100000
	return strconv.Itoa(num)
}
