package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"xml-ussd-client/client"
	"xml-ussd-client/models"
	"xml-ussd-client/utils"
)

func SendUSSDRequest(ctx context.Context, client *client.Client, requestID, msisdn, starCode, clientID string) (*models.USSDResponse, error) {
	req := models.USSDRequest{
		RequestID:    requestID,
		MSISDN:       msisdn,
		StarCode:     starCode,
		ClientID:     clientID,
		Phase:        "2",
		DCS:          "15",
		MsgType:      "4",
		UserData:     "USSD Test Message",
		EndOfSession: "0",
	}

	xmlData, err := xml.MarshalIndent(req, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to create USSD request: %v", err)
	}

	utils.LogInfo("Sending USSD request for MSISDN: %s, StarCode: %s", msisdn, starCode)

	if err := client.Write(append(xmlData, '\n')); err != nil {
		return nil, fmt.Errorf("failed to send USSD request: %v", err)
	}

	response := make([]byte, 4096)
	n, err := client.Read(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read USSD response: %v", err)
	}

	var ussdResponse models.USSDResponse
	if err := xml.Unmarshal(response[:n], &ussdResponse); err != nil {
		return nil, fmt.Errorf("failed to parse USSD response: %v", err)
	}

	if !ussdResponse.IsSuccess() {
		utils.LogError("USSD request failed: %s", ussdResponse.GetMessage())
		return &ussdResponse, fmt.Errorf("USSD request failed: %s", ussdResponse.GetMessage())
	}

	utils.LogInfo("USSD request successful")
	return &ussdResponse, nil
}
