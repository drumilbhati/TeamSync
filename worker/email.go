package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/drumilbhati/teamsync/logs"
	"github.com/drumilbhati/teamsync/utils"
	"github.com/hibiken/asynq"
)

// Unique name for task type
const TypeEmailDelivery = "email:deliver"

type EmailDeliveryPayload struct {
	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`
	OTP       string `json:"otp"`
}

/*	Producer Logic (Used by controller)	 */

// New NewEmailDeliveryTask creates a task to be queued in Redis
func NewEmailDeliveryTask(userEmail, userName, otp string) (*asynq.Task, error) {
	payload := EmailDeliveryPayload{
		UserEmail: userEmail,
		UserName:  userName,
		OTP:       otp,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Return a new task with Type and Payload
	return asynq.NewTask(TypeEmailDelivery, payloadBytes), nil
}

/*	Consumer Logic (Used by Background Worker) */

// HandleEmailDeliveryTask us the code that actually runs in the background
func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed%v: %w", err, asynq.SkipRetry)
	}

	logs.Log.Infof("Sending email to User: %s", p.UserEmail)

	err := utils.SendOTP(p.UserEmail, p.UserEmail, p.OTP)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	logs.Log.Infof("Email sent successfully to: %s", p.UserEmail)
	return nil
}
