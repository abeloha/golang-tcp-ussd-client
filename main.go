package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xml-ussd-client/client"
	"xml-ussd-client/config"
	"xml-ussd-client/handlers"
	"xml-ussd-client/utils"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		utils.LogInfo("Received shutdown signal")
		cancel()
	}()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize client with timeouts
	tcpClient := client.NewClient(
		60*time.Second, // dial timeout
		60*time.Second, // read timeout
		60*time.Second, // write timeout
	)

	// Connect to server
	if err := tcpClient.Connect(ctx, cfg.ServerHost, cfg.ServerPort); err != nil {
		utils.LogError("Failed to connect to server: %v", err)
		os.Exit(1)
	}
	defer tcpClient.Close()

	utils.LogInfo("Connected to server %s:%s", cfg.ServerHost, cfg.ServerPort)

	// Example: Login Request
	loginResp, err := handlers.SendLoginRequest(ctx, tcpClient, cfg.Username, cfg.Password)
	if err != nil {
		utils.LogError("Login failed: %v", err)
		os.Exit(1)
	}
	utils.LogInfo("Login successful: %s", loginResp.Message)

	// Example: Send USSD Request
	ussdResp, err := handlers.SendUSSDRequest(ctx, tcpClient, "12345", "2348123456789", "*123#", cfg.ClientID)
	if err != nil {
		utils.LogError("USSD request failed: %v", err)
		os.Exit(1)
	}
	utils.LogInfo("USSD request successful: %s", ussdResp.Message)

	// Example: Enquire Link
	enqResp, err := handlers.SendEnquireLink(ctx, tcpClient)
	if err != nil {
		utils.LogError("Enquire link failed: %v", err)
		os.Exit(1)
	}
	utils.LogInfo("Enquire link successful: %s", enqResp.Status)

	// Wait for shutdown signal or context cancellation
	<-ctx.Done()
	utils.LogInfo("Shutting down client")
}
