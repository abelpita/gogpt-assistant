package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	// ChatGPT API endpoint
	chatGPTAPIURL = "https://api.openai.com/v1/chat/completions"
	// Image generation API endpoint
	imageGenerationAPIURL = "https://api.openai.com/v1/images/generations"
)

func sendPromptToChatGPT(prompt string, isImagePrompt bool) (string, error) {
	// Code to send prompt to ChatGPT API will be added here
	apiToken := os.Getenv("OPENAI_GOGPT_ASSISTANT_API_KEY")
	if apiToken == "" {
		return "", fmt.Errorf("OPENAI_GOGPT_ASSISTANT_API_KEY environment variable not set")
	}

	type OpenAISettings struct {
		URL   string `json:"url"`
		Model string `json:"model"`
		N     int    `json:"n"`
		Size  string `json:"size"`
	}

	var openAISettings OpenAISettings
	var requestBody []byte
	var err error

	if isImagePrompt {
		openAISettings = OpenAISettings{
			URL:   imageGenerationAPIURL,
			Model: "DALL·E",
			N:     1,
			Size:  "256x256",
		}

		requestBody, err = json.Marshal(map[string]interface{}{
			"prompt":          prompt,
			"n":               openAISettings.N,
			"size":            openAISettings.Size,
			"response_format": "url",
		})
	} else {
		openAISettings = OpenAISettings{
			URL:   chatGPTAPIURL,
			Model: "gpt-3.5-turbo",
			N:     1,
			Size:  "",
		}

		requestBody, err = json.Marshal(map[string]interface{}{
			"model": openAISettings.Model,
			"messages": []map[string]interface{}{
				{
					"role":    "user",
					"content": prompt,
				},
			},
			"temperature": 0.7,
		})
	}
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openAISettings.URL, bytes.NewBuffer(requestBody))
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

	type ImageGenerationResponse struct {
		Created int `json:"created"`
		Data    []struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	type ErrorResponse struct {
		Error struct {
			Message string      `json:"message"`
			Type    interface{} `json:"type"`  // Replace any with interface{}
			Param   interface{} `json:"param"` // Replace any with interface{}
		} `json:"error"`
	}

	if isImagePrompt {
		var imageGenerationResponse ImageGenerationResponse
		err = json.Unmarshal(responseBody, &imageGenerationResponse)
		if err != nil {
			return "", err
		}

		if len(imageGenerationResponse.Data) == 0 {
			var errorResponse ErrorResponse
			err = json.Unmarshal(responseBody, &errorResponse)
			if err != nil {
				return "", err
			}
			return "", fmt.Errorf(errorResponse.Error.Message)
		}

		generatedImageURL := imageGenerationResponse.Data[0].URL
		return generatedImageURL, nil
	} else {
		var chatGPTResponse ChatGPTResponse
		err = json.Unmarshal(responseBody, &chatGPTResponse)
		if err != nil {
			return "", err
		}

		generatedText := chatGPTResponse.Choices[0].Message.Content
		return generatedText, nil
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("GoGPT-Assistant")
	myWindow.Resize(fyne.NewSize(800, 600))

	promptInput := widget.NewEntry()
	promptInput.SetPlaceHolder("Enter your text prompt or image description here")

	responseLabel := widget.NewLabel("Response will be displayed here")
	isImagePromptCheckbox := widget.NewCheck("Generate Image", nil)
	isImagePromptLabel := widget.NewLabel("I want to generate an image")
	responseImage := canvas.NewImageFromFile("")
	responseImage.FillMode = canvas.ImageFillContain
	responseImage.SetMinSize(fyne.NewSize(256, 256))

	sendButton := widget.NewButton("Send", func() {
		prompt := promptInput.Text
		isImagePrompt := isImagePromptCheckbox.Checked
		response, err := sendPromptToChatGPT(prompt, isImagePrompt)
		if err != nil {
			responseLabel.SetText(fmt.Sprintf("Error: %s", err))
		} else {
			if isImagePrompt {
				resource, err := fyne.LoadResourceFromURLString(response)
				if err != nil {
					responseLabel.SetText(fmt.Sprintf("Error: %s", err))
					return
				}
				responseImage.Resource = resource
				responseImage.Refresh()
				responseLabel.SetText("")
			} else {
				responseLabel.SetText(response)
			}
		}
	})

	content := container.NewVBox(
		promptInput,
		isImagePromptCheckbox,
		isImagePromptLabel,
		sendButton,
		responseLabel,
		responseImage,
	)

	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
