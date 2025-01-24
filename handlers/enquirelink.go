package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"xml-ussd-client/client"
	"xml-ussd-client/models"
	"xml-ussd-client/utils"
)

func SendEnquireLink(ctx context.Context, client *client.Client) (*models.ENQResponse, error) {
	req := models.ENQRequest{}
	xmlData, err := xml.MarshalIndent(req, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to create enquire link request: %v", err)
	}

	utils.LogInfo("Sending enquire link request")

	if err := client.Write(append(xmlData, '\n')); err != nil {
		return nil, fmt.Errorf("failed to send enquire link request: %v", err)
	}

	response := make([]byte, 4096)
	n, err := client.Read(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read enquire link response: %v", err)
	}

	var enqResponse models.ENQResponse
	if err := xml.Unmarshal(response[:n], &enqResponse); err != nil {
		return nil, fmt.Errorf("failed to parse enquire link response: %v", err)
	}

	if !enqResponse.IsSuccess() {
		utils.LogError("Enquire link failed: %s", enqResponse.GetMessage())
		return &enqResponse, fmt.Errorf("enquire link failed: %s", enqResponse.GetMessage())
	}

	utils.LogInfo("Enquire link successful")
	return &enqResponse, nil
}
