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

		req := models.AuthRequest{
			RequestID:  time.Now().Format("20060102150405"),
			Username: username,
			Password: password,
			ApplicationID: cfg.ClientID,
		}
		
		// Log request details with masked password
		maskedPassword := "****" + password[len(password)-4:]
		utils.LogInfo("Login Request Details:")
		utils.LogInfo("Username: %s", username)
		utils.LogInfo("Password: %s", maskedPassword)

		xmlData, err := xml.MarshalIndent(req, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to create login request: %v", err)
		}

		// Log raw XML payload
		utils.LogInfo("XML Payload:\n%s", string(xmlData))

		utils.LogInfo("Sending login request for user: %s (Attempt %d)", username, attempt+1)
		
		fullPayload := append(xmlData, '\n')
		utils.LogInfo("Full Payload (with newline): %v", fullPayload)

		if err := client.Write(fullPayload); err != nil {
			lastErr = fmt.Errorf("failed to send login request on attempt %d: %v", attempt+1, err)
			utils.LogError("%v", lastErr)
			time.Sleep(2 * time.Second)
			continue
		}

		response := make([]byte, 4096)
		n, err := client.Read(response)
		if err != nil {
			lastErr = fmt.Errorf("failed to read login response on attempt %d: %v", attempt+1, err)
			utils.LogError("%v", lastErr)
			time.Sleep(2 * time.Second)
			continue
		}

		// Log raw response
		utils.LogInfo("Raw Response (bytes): %v", response[:n])
		utils.LogInfo("Raw Response (string): %s", string(response[:n]))

		var authResponse models.AuthResponse
		if err := xml.Unmarshal(response[:n], &authResponse); err != nil {
			utils.LogError("XML Unmarshal Error: %v", err)
			return nil, fmt.Errorf("failed to parse login response: %v", err)
		}

		// Log parsed response details
		utils.LogInfo("Auth Response Details:")
		utils.LogInfo("Success: %v", authResponse.IsSuccess())
		utils.LogInfo("Message: %s", authResponse.GetMessage())

		if !authResponse.IsSuccess() {
			utils.LogError("Login failed: %s", authResponse.GetMessage())
			return &authResponse, fmt.Errorf("login failed: %s", authResponse.GetMessage())
		}

		utils.LogInfo("Login successful")
		return &authResponse, nil
	}

	return nil, fmt.Errorf("login failed after %d attempts: %v", maxRetries, lastErr)
}
