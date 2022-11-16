package eda

const (
	AppointmentRequestChannel  = "appointment-req-channel"
	AppointmentResponseChannel = "appointment-resp-channel"

	/** Telehealers RESP backend constants **/
	restBackendAddress = "http://127.0.0.1:1234"
	restAuthTokenName  = "telehealers-token"
	restAuthPass       = "letmein"
)

var authHeader = map[string][]string{restAuthTokenName: []string{restAuthPass}}
