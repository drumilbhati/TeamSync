package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/drumilbhati/teamsync/models"
	"google.golang.org/genai"
)

func EnhanceTask(task *models.Task) (*models.Task, error) {
	// 1. Convert input task to JSON for the prompt
	inputBytes, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input task: %v", err)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %v", err)
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
4. **Structure:** Return the EXACT same JSON structure (fields and types) as the input. Do not add or remove fields. Only modify the values of 'title', 'description', 'status', and 'priority' if needed. Preserve all IDs and dates.

Return ONLY the raw JSON string. No markdown formatting.`, string(inputBytes))

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("ai generation failed: %v", err)
	}

	// 2. Clean the response (LLMs sometimes add ```json ... ``` blocks even when told not to)
	responseText := result.Text()
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	// 3. Unmarshal the JSON back into a Task struct
	var enhancedTask models.Task
	err = json.Unmarshal([]byte(responseText), &enhancedTask)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v. Response was: %s", err, responseText)
	}

	return &enhancedTask, nil
}
