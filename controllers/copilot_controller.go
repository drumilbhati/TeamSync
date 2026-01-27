package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/genai"
)

type TaskInput struct {
	Title       string `json:"title"`
	Description struct {
		String string `json:"String"`
		Valid  bool   `json:"Valid"`
	} `json:"description"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	TeamID   int    `json:"team_id"`
	UserID   int    `json:"user_id"`
}

func Describe(w http.ResponseWriter, r *http.Request) {
	var input TaskInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prompt := fmt.Sprintf(`The data given to you is related to a task that is being generated on a task management platform. The data may be vague, incomplete and not useful to the team of the user. 
			Give a detailed description for this task based on the given input data in order to make this more useful and descriptive. 
			
			Task Details:
			Title: %s
			Description: %s
			Status: %s
			Priority: %s
			
			Give two fields in your response object:
			1. The new and improved title for the task
			2. The new and improve description for the task
			Only give these two fields, no other text is expected from you, neither should you give.
			The fields of output should be same as the fields of input, just modify them to be more descriptive.
			Match the case with that of the input fields (everything in fieldname is lowercase.
			Give a well formatted json object in return with all the five fields as given to you.`,
		input.Title, input.Description.String, input.Status, input.Priority)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(result.Text())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Text())
}
