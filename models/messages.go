package models

import "encoding/xml"

// Login Request and Response
type AuthRequest struct {
	XMLName  xml.Name `xml:"AUTHRequest"`
	Username string   `xml:"username"`
	Password string   `xml:"password"`
}

type AuthResponse struct {
	XMLName xml.Name `xml:"AUTHResponse"`
	Status  string   `xml:"status"`
	Message string   `xml:"message"`
}

// USSD Request and Response
type USSDRequest struct {
	XMLName      xml.Name `xml:"USSDRequest"`
	RequestID    string   `xml:"requestId"`
	MSISDN       string   `xml:"msisdn"`
	StarCode     string   `xml:"starCode"`
	ClientID     string   `xml:"clientId"`
	Phase        string   `xml:"phase"`
	DCS          string   `xml:"dcs"`
	MsgType      string   `xml:"msgtype"`
	UserData     string   `xml:"userdata"`
	EndOfSession string   `xml:"EndofSession"`
}

type USSDResponse struct {
	XMLName      xml.Name `xml:"USSDResponse"`
	RequestID    string   `xml:"requestId"`
	Status       string   `xml:"status"`
	Message      string   `xml:"message"`
	UserData     string   `xml:"userdata,omitempty"`
	EndOfSession string   `xml:"EndofSession,omitempty"`
}

// Enquire Link Request and Response
type ENQRequest struct {
	XMLName xml.Name `xml:"ENQRequest"`
}

type ENQResponse struct {
	XMLName xml.Name `xml:"ENQResponse"`
	Status  string   `xml:"status"`
}

// Error types
const (
	StatusSuccess = "0"
	StatusError   = "1"
)

// Response interface for common response handling
type Response interface {
	IsSuccess() bool
	GetMessage() string
}

// Implementation of Response interface for each response type
func (r *AuthResponse) IsSuccess() bool {
	return r.Status == StatusSuccess
}

func (r *AuthResponse) GetMessage() string {
	return r.Message
}

func (r *USSDResponse) IsSuccess() bool {
	return r.Status == StatusSuccess
}

func (r *USSDResponse) GetMessage() string {
	return r.Message
}

func (r *ENQResponse) IsSuccess() bool {
	return r.Status == StatusSuccess
}

func (r *ENQResponse) GetMessage() string {
	return r.Status
}
