package envelope

import (
	"encoding/json"
	"fmt"
	"git.wisehodl.dev/jay/go-roots-ws/errors"
)

// CheckArrayLength is a helper function that ensures the JSON array has at
// least the minimum length required
func CheckArrayLength(arr []json.RawMessage, minLen int) error {
	if len(arr) < minLen {
		return fmt.Errorf("%w: expected %d elements, got %d", errors.InvalidEnvelope, minLen, len(arr))
	}
	return nil
}

// CheckLabel is a helper function that verifies that the envelope label
// matches the expected one
func CheckLabel(got, want string) error {
	if got != want {
		return fmt.Errorf("%w: expected %s, got %s", errors.WrongEnvelopeLabel, want, got)
	}
	return nil
}

// ParseElement is a helper function that unmarshals an array element into the
// provided value
func ParseElement(element json.RawMessage, value interface{}, position string) error {
	if err := json.Unmarshal(element, value); err != nil {
		return fmt.Errorf("%w: %s is not the expected type", errors.WrongFieldType, position)
	}
	return nil
}

// FindEvent extracts an event from an EVENT envelope
// Expected Format: ["EVENT", event]
func FindEvent(env Envelope) ([]byte, error) {
	var arr []json.RawMessage
	if err := json.Unmarshal(env, &arr); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err := CheckArrayLength(arr, 2); err != nil {
		return nil, err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return nil, err
	}

	if err := CheckLabel(label, "EVENT"); err != nil {
		return nil, err
	}

	return arr[1], nil
}

// FindEventWithReq extracts an event from an EVENT envelope with a subscription ID.
// Expected Format: ["EVENT", "SUBID", event]
func FindEventWithReq(env Envelope) (string, []byte, error) {
	var arr []json.RawMessage
	if err := json.Unmarshal(env, &arr); err != nil {
		return "", nil, fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err := CheckArrayLength(arr, 3); err != nil {
		return "", nil, err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return "", nil, err
	}

	if err := CheckLabel(label, "EVENT"); err != nil {
		return "", nil, err
	}

	var req string
	if err := ParseElement(arr[1], &req, "request id"); err != nil {
		return "", nil, err
	}

	return req, arr[2], nil
}

// FindSubscriptionEvent extracts an event and subscription ID from an EVENT envelope.
// Expected Format: ["EVENT", subID, event]
func FindSubscriptionEvent(env Envelope) (subID string, event []byte, err error) {
	var arr []json.RawMessage
	if err = json.Unmarshal(env, &arr); err != nil {
		return "", nil, fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err = CheckArrayLength(arr, 3); err != nil {
		return "", nil, err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return "", nil, err
	}

	if err = CheckLabel(label, "EVENT"); err != nil {
		return "", nil, err
	}

	if err = ParseElement(arr[1], &subID, "subscription ID"); err != nil {
		return "", nil, err
	}

	return subID, arr[2], nil
}

// FindOK extracts eventID, status, and message from an OK envelope.
// Expected Format: ["OK", eventID, status, message]
func FindOK(env Envelope) (eventID string, status bool, message string, err error) {
	var arr []json.RawMessage
	if err = json.Unmarshal(env, &arr); err != nil {
		return "", false, "", fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err = CheckArrayLength(arr, 4); err != nil {
		return "", false, "", err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return "", false, "", err
	}

	if err = CheckLabel(label, "OK"); err != nil {
		return "", false, "", err
	}

	if err = ParseElement(arr[1], &eventID, "event ID"); err != nil {
		return "", false, "", err
	}

	if err = ParseElement(arr[2], &status, "status"); err != nil {
		return "", false, "", err
	}

	if err = ParseElement(arr[3], &message, "message"); err != nil {
		return "", false, "", err
	}

	return eventID, status, message, nil
}

// FindReq extracts subscription ID and filters from a REQ envelope.
// Expected Format: ["REQ", subID, filter1, filter2, ...]
func FindReq(env Envelope) (subID string, filters [][]byte, err error) {
	var arr []json.RawMessage
	if err = json.Unmarshal(env, &arr); err != nil {
		return "", nil, fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err = CheckArrayLength(arr, 2); err != nil {
		return "", nil, err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return "", nil, err
	}

	if err = CheckLabel(label, "REQ"); err != nil {
		return "", nil, err
	}

	if err = ParseElement(arr[1], &subID, "subscription ID"); err != nil {
		return "", nil, err
	}

	filters = make([][]byte, 0, len(arr)-2)
	for i := 2; i < len(arr); i++ {
		filters = append(filters, arr[i])
	}

	return subID, filters, nil
}

// FindEOSE extracts subscription ID from an EOSE envelope.
// Expected Format: ["EOSE", subID]
func FindEOSE(env Envelope) (subID string, err error) {
	var arr []json.RawMessage
	if err = json.Unmarshal(env, &arr); err != nil {
		return "", fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err = CheckArrayLength(arr, 2); err != nil {
		return "", err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return "", err
	}

	if err = CheckLabel(label, "EOSE"); err != nil {
		return "", err
	}

	if err = ParseElement(arr[1], &subID, "subscription ID"); err != nil {
		return "", err
	}

	return subID, nil
}

// FindClose extracts subscription ID from a CLOSE envelope.
// Expected Format: ["CLOSE", subID]
func FindClose(env Envelope) (subID string, err error) {
	var arr []json.RawMessage
	if err = json.Unmarshal(env, &arr); err != nil {
		return "", fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err = CheckArrayLength(arr, 2); err != nil {
		return "", err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return "", err
	}

	if err = CheckLabel(label, "CLOSE"); err != nil {
		return "", err
	}

	if err = ParseElement(arr[1], &subID, "subscription ID"); err != nil {
		return "", err
	}

	return subID, nil
}

// FindClosed extracts subscription ID and message from a CLOSED envelope.
// Expected Format: ["CLOSED", subID, message]
func FindClosed(env Envelope) (subID string, message string, err error) {
	var arr []json.RawMessage
	if err = json.Unmarshal(env, &arr); err != nil {
		return "", "", fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err = CheckArrayLength(arr, 3); err != nil {
		return "", "", err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return "", "", err
	}

	if err = CheckLabel(label, "CLOSED"); err != nil {
		return "", "", err
	}

	if err = ParseElement(arr[1], &subID, "subscription ID"); err != nil {
		return "", "", err
	}

	if err = ParseElement(arr[2], &message, "message"); err != nil {
		return "", "", err
	}

	return subID, message, nil
}

// FindNotice extracts message from a NOTICE envelope.
// Expected Format: ["NOTICE", message]
func FindNotice(env Envelope) (message string, err error) {
	var arr []json.RawMessage
	if err = json.Unmarshal(env, &arr); err != nil {
		return "", fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err = CheckArrayLength(arr, 2); err != nil {
		return "", err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return "", err
	}

	if err = CheckLabel(label, "NOTICE"); err != nil {
		return "", err
	}

	if err = ParseElement(arr[1], &message, "message"); err != nil {
		return "", err
	}

	return message, nil
}

// FindAuthChallenge extracts challenge from an AUTH challenge envelope.
// Expected Format: ["AUTH", challenge]
func FindAuthChallenge(env Envelope) (challenge string, err error) {
	var arr []json.RawMessage
	if err = json.Unmarshal(env, &arr); err != nil {
		return "", fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err = CheckArrayLength(arr, 2); err != nil {
		return "", err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return "", err
	}

	if err = CheckLabel(label, "AUTH"); err != nil {
		return "", err
	}

	// Check if the second element is a string (AUTH challenge)
	if err = ParseElement(arr[1], &challenge, "challenge"); err != nil {
		return "", err
	}

	return challenge, nil
}

// FindAuthResponse extracts event from an AUTH response envelope.
// Expected Format: ["AUTH", event]
func FindAuthResponse(env Envelope) (event []byte, err error) {
	var arr []json.RawMessage
	if err = json.Unmarshal(env, &arr); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if err = CheckArrayLength(arr, 2); err != nil {
		return nil, err
	}

	var label string
	if err := ParseElement(arr[0], &label, "envelope label"); err != nil {
		return nil, err
	}

	if err = CheckLabel(label, "AUTH"); err != nil {
		return nil, err
	}

	return arr[1], nil
}
