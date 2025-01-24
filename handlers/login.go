package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"xml-ussd-client/client"
	"xml-ussd-client/models"
	"xml-ussd-client/utils"
)

func SendLoginRequest(ctx context.Context, client *client.Client, username, password string) (*models.AuthResponse, error) {
	req := models.AuthRequest{
		Username: username,
		Password: password,
	}
	
	xmlData, err := xml.MarshalIndent(req, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %v", err)
	}

	utils.LogInfo("Sending login request for user: %s", username)
	
	if err := client.Write(append(xmlData, '\n')); err != nil {
		return nil, fmt.Errorf("failed to send login request: %v", err)
	}

	response := make([]byte, 4096)
	n, err := client.Read(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read login response: %v", err)
	}

	var authResponse models.AuthResponse
	if err := xml.Unmarshal(response[:n], &authResponse); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %v", err)
	}

	if !authResponse.IsSuccess() {
		utils.LogError("Login failed: %s", authResponse.GetMessage())
		return &authResponse, fmt.Errorf("login failed: %s", authResponse.GetMessage())
	}

	utils.LogInfo("Login successful")
	return &authResponse, nil
}
