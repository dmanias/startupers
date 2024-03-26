package ai

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

// AskAI adjusted to use the go-openai package
func AskAI(apiKey string, prompt string) (string, error) {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err // Return an empty string along with the error
	}

	// Ensure there is at least one choice and its message content is accessible before accessing it
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		fmt.Println(resp.Choices[0].Message.Content)
		return resp.Choices[0].Message.Content, nil
	}

	// If no choices are available or the message content is empty, return an appropriate error
	return "", fmt.Errorf("no response received from the AI")
}
