// Package errors defines standard error types used throughout the roots-ws library.
package errors

import (
	"errors"
)

var (
	// Data Structure Errors

	// InvalidJSON indicates that a byte sequence could not be parsed as valid JSON.
	// This is typically returned when unmarshaling fails during envelope processing.
	InvalidJSON = errors.New("invalid JSON")

	// MissingField indicates that a required field is absent from a data structure.
	// This is returned when validating that all mandatory components are present.
	MissingField = errors.New("missing required field")

	// WrongFieldType indicates that a field's type does not match the expected type.
	// This is returned when unmarshaling a specific value fails due to type mismatch.
	WrongFieldType = errors.New("wrong field type")

	// Envelope Errors

	// InvalidEnvelope indicates that a message does not conform to the Nostr envelope structure.
	// This typically occurs when an array has incorrect number of elements for its message type.
	InvalidEnvelope = errors.New("invalid envelope format")

	// WrongEnvelopeLabel indicates that an envelope's label does not match the expected type.
	// This is returned when attempting to parse an envelope using a Find function that
	// expects a different label than what was provided.
	WrongEnvelopeLabel = errors.New("wrong envelope label")
)
