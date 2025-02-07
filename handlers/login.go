package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"
	"xml-ussd-client/client"
	"xml-ussd-client/config"
	"xml-ussd-client/models"
	"xml-ussd-client/utils"
)

func SendLoginRequest(ctx context.Context, client *client.Client, cfg config.Config) (*models.AuthResponse, error) {

	username := cfg.Username
	password := cfg.Password

	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Log connection details
		utils.LogInfo("Login Attempt %d: Connecting to server %s:%s", attempt+1, cfg.ServerHost, cfg.ServerPort)

		// Reconnect on each retry to handle potential connection issues
		if err := client.Close(); err != nil {
			utils.LogError("Error closing previous connection: %v", err)
		}

		// Recreate connection for each attempt
		if err := client.Connect(ctx, cfg.ServerHost, cfg.ServerPort); err != nil {
			lastErr = fmt.Errorf("failed to reconnect on attempt %d: %v", attempt+1, err)
			utils.LogError("%v", lastErr)
			time.Sleep(2 * time.Second) // Wait before retrying
			continue
		}

		// Generate unique session ID for this request
		sessionID := fmt.Sprintf("%x", time.Now().UnixNano())

		req := models.AuthRequest{
			RequestID:     sessionID,
			Username:       username,
			Password:       password,
			ApplicationID:  cfg.ClientID,
		}
		
		// Marshal XML payload
		xmlData, err := xml.MarshalIndent(req, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to create login request: %v", err)
		}

		// Add newline to XML payload
		fullPayload := append(xmlData, '\n')
		
		// Write message with header using generated sessionID
		err = client.WriteMessageWithHeader(sessionID, "LOGIN", fullPayload)
		if err != nil {
			lastErr = fmt.Errorf("failed to send login request on attempt %d: %v", attempt+1, err)
			utils.LogError("%v", lastErr)
			time.Sleep(2 * time.Second)
			continue
		}

		// Read response with header
		receivedSessionID, responsePayload, err := client.ReadMessageWithHeader()
		if err != nil {
			lastErr = fmt.Errorf("failed to read login response on attempt %d: %v", attempt+1, err)
			utils.LogError("%v", lastErr)
			time.Sleep(2 * time.Second)
			continue
		}

		// Log session ID for correlation
		utils.LogInfo("Sent Session ID: %s", sessionID)
		utils.LogInfo("Received Session ID: %s", receivedSessionID)

		var authResponse models.AuthResponse
		if err := xml.Unmarshal(responsePayload, &authResponse); err != nil {
			utils.LogError("XML Unmarshal Error: %v", err)
			return nil, fmt.Errorf("failed to parse login response: %v", err)
		}

		if !authResponse.IsSuccess() {
			utils.LogError("Login failed: %s", authResponse.GetMessage())
			return &authResponse, fmt.Errorf("login failed: %s", authResponse.GetMessage())
		}

		utils.LogInfo("Login successful")
		return &authResponse, nil
	}

	return nil, fmt.Errorf("login failed after %d attempts: %v", maxRetries, lastErr)
}
