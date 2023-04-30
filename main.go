package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("GoGPT-Assistant")

	// UI components will be added here
	// Create inptu field for prompts
	promptInput := widget.NewEntry()
	promptInput.SetPlaceHolder("Enter your prompt here")

	// Create area to display response
	responseLabel := widget.NewLabel("Response will be displayed here")

	// Create button send button and its action
	sendButton := widget.NewButton("Send", func() {
		// Code to send prompt to ChatGPT API and display response will be added here
	})

	// Organize layout using containers
	content := container.NewVBox(
		promptInput,
		sendButton,
		responseLabel,
	)

	myWindow.SetContent(content)

	myWindow.ShowAndRun()
}
