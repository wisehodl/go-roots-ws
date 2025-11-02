package envelope

import (
	"bytes"
	"strconv"
)

// EncloseEvent creates an EVENT envelope for publishing events.
// It wraps the provided event JSON in the format ["EVENT", event].
func EncloseEvent(event []byte) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["EVENT",`)
	buf.Write(event)
	buf.WriteByte(']')
	return buf.Bytes()
}

// EncloseOK creates an OK envelope acknowledging receipt of an event.
// Format: ["OK", eventID, status, message]
func EncloseOK(eventID string, status bool, message string) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["OK","`)
	buf.WriteString(eventID)
	buf.WriteString(`",`)
	buf.WriteString(strconv.FormatBool(status))
	buf.WriteString(`,"`)
	buf.WriteString(message)
	buf.WriteString(`"]`)
	return buf.Bytes()
}

// EncloseReq creates a REQ envelope for subscription requests.
// Format: ["REQ", subID, filter1, filter2, ...]
func EncloseReq(subID string, filters [][]byte) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["REQ","`)
	buf.WriteString(subID)
	buf.WriteString(`"`)

	for _, filter := range filters {
		buf.WriteString(`,`)
		buf.Write(filter)
	}

	buf.WriteByte(']')
	return buf.Bytes()
}

// EncloseSubscriptionEvent creates an EVENT envelope for delivering subscription events.
// Format: ["EVENT", subID, event]
func EncloseSubscriptionEvent(subID string, event []byte) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["EVENT","`)
	buf.WriteString(subID)
	buf.WriteString(`",`)
	buf.Write(event)
	buf.WriteByte(']')
	return buf.Bytes()
}

// EncloseEOSE creates an EOSE (End of Stored Events) envelope.
// Format: ["EOSE", subID]
func EncloseEOSE(subID string) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["EOSE","`)
	buf.WriteString(subID)
	buf.WriteString(`"]`)
	return buf.Bytes()
}

// EncloseClose creates a CLOSE envelope for ending a subscription.
// Format: ["CLOSE", subID]
func EncloseClose(subID string) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["CLOSE","`)
	buf.WriteString(subID)
	buf.WriteString(`"]`)
	return buf.Bytes()
}

// EncloseClosed creates a CLOSED envelope for indicating a terminated subscription.
// Format: ["CLOSED", subID, message]
func EncloseClosed(subID string, message string) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["CLOSED","`)
	buf.WriteString(subID)
	buf.WriteString(`","`)
	buf.WriteString(message)
	buf.WriteString(`"]`)
	return buf.Bytes()
}

// EncloseNotice creates a NOTICE envelope for responder messages.
// Format: ["NOTICE", message]
func EncloseNotice(message string) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["NOTICE","`)
	buf.WriteString(message)
	buf.WriteString(`"]`)
	return buf.Bytes()
}

// EncloseAuthChallenge creates an AUTH challenge envelope.
// Format: ["AUTH", challenge]
func EncloseAuthChallenge(challenge string) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["AUTH","`)
	buf.WriteString(challenge)
	buf.WriteString(`"]`)
	return buf.Bytes()
}

// EncloseAuthResponse creates an AUTH response envelope.
// Format: ["AUTH", event]
func EncloseAuthResponse(event []byte) Envelope {
	var buf bytes.Buffer
	buf.WriteString(`["AUTH",`)
	buf.Write(event)
	buf.WriteByte(']')
	return buf.Bytes()
}
