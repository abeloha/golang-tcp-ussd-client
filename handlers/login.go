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
			Username: username,
			Password: password,
		}
		
		xmlData, err := xml.MarshalIndent(req, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to create login request: %v", err)
		}

		utils.LogInfo("Sending login request for user: %s (Attempt %d)", username, attempt+1)
		
		if err := client.Write(append(xmlData, '\n')); err != nil {
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

	return nil, fmt.Errorf("login failed after %d attempts: %v", maxRetries, lastErr)
}
