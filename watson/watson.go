package watson

import (
	"fmt"

	"github.com/IBM/go-sdk-core/core"
	"github.com/watson-developer-cloud/go-sdk/assistantv2"
)

func Send(text string) (*assistantv2.MessageResponse, error) {
	apiVersion := "2018-09-20"
	apiURL := "https://gateway-lon.watsonplatform.net/assistant/api"
	apiUsername := "apikey"
	apiPassword := "EsYpIj2JZKGHDkxe4V2yQV7fOfM-Py4bTfC4q0h6nRn6"
	assistantID := "fea2548f-d8fb-4767-a610-f73810142612"
	//intentName := "General_Events"
	//intentConfidence := 1.0
	service, serviceErr := assistantv2.NewAssistantV2(&assistantv2.AssistantV2Options{
		Version:  apiVersion,
		URL:      apiURL,
		Username: apiUsername,
		Password: apiPassword,
	})
	if serviceErr != nil {
		return "", fmt.Errorf("error creating service: %v", serviceErr)
	}
	createSessionResponse, createSessionResponseErr := service.CreateSession(service.NewCreateSessionOptions(assistantID))
	if createSessionResponseErr != nil {
		return "", fmt.Errorf("error creating service: %v", createSessionResponseErr)
	}
	createSessionResult := service.GetCreateSessionResult(createSessionResponse)
	sessionID := *createSessionResult.SessionID
	fmt.Print(sessionID, "\n")
	messageResponse, messageResponseErr := service.Message(
		&assistantv2.MessageOptions{
			AssistantID: core.StringPtr(assistantID),
			SessionID:   core.StringPtr(sessionID),
			Input: &assistantv2.MessageInput{
				MessageType: core.StringPtr("text"),
				Text:        core.StringPtr(text),
				//Intents: []assistantv2.RuntimeIntent{{Intent: &intentName, Confidence: &intentConfidence}},
			},
		},
	)
	if messageResponseErr != nil {
		return "", fmt.Errorf("error fetching response: %v", messageResponseErr)
	}
	messageResult := service.GetMessageResult(messageResponse)
	return messageResult, nil
}
