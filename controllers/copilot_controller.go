package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/drumilbhati/teamsync/models"
	"google.golang.org/genai"
)

func Describe(w http.ResponseWriter, r *http.Request) {
	var input models.Task
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inputBytes, err := json.Marshal(input)
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

	prompt := fmt.Sprintf(`You are an AI assistant for a task management platform.
You will receive a JSON object representing a task.
Your goal is to enhance the task details to be more clear, professional, and actionable.

Input Task (JSON):
%s

Enhancement Rules:
1. **Title:** Make it concise but descriptive.
2. **Description:** Expand on the description to provide context, potential steps, or necessary details based on the title and existing description. Ensure it is well-formatted.
3. **Consistency:** Ensure the 'status' and 'priority' match the context of the task. If the text implies urgency, ensure priority is 'high'.
4. **Structure:** Return the EXACT same JSON structure (fields and types) as the input. Do not add or remove fields. Only modify the values of 'title', 'description' (including 'String' and 'Valid' fields inside it), 'status', and 'priority' if needed. Preserve all IDs and dates.

Return ONLY the raw JSON string. No markdown formatting.`, string(inputBytes))

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

	w.Header().Set("Content-Type", "application/json")
	// The result.Text() should be the JSON string. We write it directly.
	w.Write([]byte(result.Text()))
}
