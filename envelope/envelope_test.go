package envelope

import (
	"git.wisehodl.dev/jay/go-roots-ws/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLabel(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantLabel   string
		wantErr     error
		wantErrText string
	}{
		{
			name:      "valid envelope with EVENT label",
			env:       []byte(`["EVENT",{"id":"abc123"}]`),
			wantLabel: "EVENT",
		},
		{
			name:      "valid envelope with custom label",
			env:       []byte(`["TEST",{"data":"value"}]`),
			wantLabel: "TEST",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "empty array",
			env:         []byte(`[]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "empty envelope",
		},
		{
			name:        "label not a string",
			env:         []byte(`[123,{"id":"abc123"}]`),
			wantErr:     errors.WrongFieldType,
			wantErrText: "label is not a string",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetLabel(tc.env)

			if tc.wantErr != nil || tc.wantErrText != "" {
				if tc.wantErr != nil {
					assert.ErrorIs(t, err, tc.wantErr)
				}
				if tc.wantErrText != "" {
					assert.ErrorContains(t, err, tc.wantErrText)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantLabel, got)
		})
	}
}

func TestGetStandardLabels(t *testing.T) {
	expected := map[string]struct{}{
		"EVENT":  {},
		"REQ":    {},
		"CLOSE":  {},
		"CLOSED": {},
		"EOSE":   {},
		"NOTICE": {},
		"OK":     {},
		"AUTH":   {},
	}

	labels := GetStandardLabels()

	// Check that we have the exact same number of labels
	assert.Equal(t, len(expected), len(labels))

	// Check that all expected labels are present
	for label := range expected {
		_, exists := labels[label]
		assert.True(t, exists, "Expected standard label %s not found", label)
	}
}

func TestIsStandardLabel(t *testing.T) {
	standardCases := []string{
		"EVENT", "REQ", "CLOSE", "CLOSED", "EOSE", "NOTICE", "OK", "AUTH",
	}

	nonStandardCases := []string{
		"TEST", "CUSTOM", "event", "REQ1", "",
	}

	for _, label := range standardCases {
		t.Run(label, func(t *testing.T) {
			assert.True(t, IsStandardLabel(label), "Label %s should be standard", label)
		})
	}

	for _, label := range nonStandardCases {
		t.Run(label, func(t *testing.T) {
			assert.False(t, IsStandardLabel(label), "Label %s should not be standard", label)
		})
	}
}
