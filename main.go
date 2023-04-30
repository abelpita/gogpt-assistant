package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	chatGPTAPIURL = "https://api.openai.com/v1/chat/completions"
)

func sendPromptToChatGPT(prompt string) (string, error) {
	// Code to send prompt to ChatGPT API will be added here
	apiToken := os.Getenv("OPENAI_GOGPT_ASSISTANT_API_KEY")
	if apiToken == "" {
		return "", fmt.Errorf("OPENAI_GOGPT_ASSISTANT_API_KEY environment variable not set")
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
	})

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", chatGPTAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type ChatGPTResponse struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Model   string `json:"model"`
		Usage   struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
			Index        int    `json:"index"`
		} `json:"choices"`
	}

	var chatGPTResponse ChatGPTResponse
	err = json.Unmarshal(responseBody, &chatGPTResponse)
	if err != nil {
		return "", err
	}

	generatedText := chatGPTResponse.Choices[0].Message.Content
	return generatedText, nil
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("GoGPT-Assistant")

	promptInput := widget.NewEntry()
	promptInput.SetPlaceHolder("Enter your prompt here")

	responseLabel := widget.NewLabel("Response will be displayed here")

	sendButton := widget.NewButton("Send", func() {
		prompt := promptInput.Text
		response, err := sendPromptToChatGPT(prompt)
		if err != nil {
			responseLabel.SetText(fmt.Sprintf("Error: %s", err))
		} else {
			responseLabel.SetText(response)
		}
	})

	content := container.NewVBox(
		promptInput,
		sendButton,
		responseLabel,
	)

	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
