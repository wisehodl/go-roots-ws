// Package envelope provides types and functions for working with Nostr protocol
// websocket messages. It defines the Envelope type representing a Nostr message
// and offers utilities for creating, parsing, and validating standardized message
// formats.
package envelope

import (
	"encoding/json"
	"fmt"
	"git.wisehodl.dev/jay/go-roots-ws/errors"
)

// Envelope represents a Nostr websocket message.
type Envelope []byte

// GetLabel extracts the message label from an envelope.
// Returns the label as a string or an error if the envelope is malformed.
func GetLabel(env Envelope) (string, error) {
	var arr []json.RawMessage
	if err := json.Unmarshal(env, &arr); err != nil {
		return "", fmt.Errorf("%w: %v", errors.InvalidJSON, err)
	}

	if len(arr) < 1 {
		return "", fmt.Errorf("%w: empty envelope", errors.InvalidEnvelope)
	}

	var label string
	if err := json.Unmarshal(arr[0], &label); err != nil {
		return "", fmt.Errorf("%w: label is not a string", errors.WrongFieldType)
	}

	return label, nil
}

// GetStandardLabels returns a set of standard Nostr websocket message labels
func GetStandardLabels() map[string]struct{} {
	return map[string]struct{}{
		"EVENT":  {},
		"REQ":    {},
		"CLOSE":  {},
		"CLOSED": {},
		"EOSE":   {},
		"NOTICE": {},
		"OK":     {},
		"AUTH":   {},
	}
}

// IsStandardLabel checks if the given label is a standard Nostr websocket message label
func IsStandardLabel(label string) bool {
	labels := GetStandardLabels()
	_, ok := labels[label]
	return ok
}
