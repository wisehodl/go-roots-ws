package envelope

import (
	"testing"

	"git.wisehodl.dev/jay/go-roots-ws/errors"
	"github.com/stretchr/testify/assert"
)

func TestFindEvent(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantEvent   []byte
		wantErr     error
		wantErrText string
	}{
		{
			name:      "valid event",
			env:       []byte(`["EVENT",{"id":"abc123","kind":1}]`),
			wantEvent: []byte(`{"id":"abc123","kind":1}`),
		},
		{
			name:        "wrong label",
			env:         []byte(`["REQ",{"id":"abc123","kind":1}]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected EVENT, got REQ",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["EVENT"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 2 elements, got 1",
		},
		{
			name:      "extraneous elements",
			env:       []byte(`["EVENT",{"id":"abc123"},"extra"]`),
			wantEvent: []byte(`{"id":"abc123"}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FindEvent(tc.env)

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
			assert.Equal(t, tc.wantEvent, got)
		})
	}
}

func TestFindEventWithReq(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantReq     string
		wantEvent   []byte
		wantErr     error
		wantErrText string
	}{
		{
			name:      "valid event",
			env:       []byte(`["EVENT","SUBID",{"id":"abc123","kind":1}]`),
			wantReq:   "SUBID",
			wantEvent: []byte(`{"id":"abc123","kind":1}`),
		},
		{
			name:        "wrong label",
			env:         []byte(`["REQ","SUBID",{"id":"abc123","kind":1}]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected EVENT, got REQ",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["EVENT","SUBID"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 3 elements, got 2",
		},
		{
			name:      "extraneous elements",
			env:       []byte(`["EVENT","SUBID",{"id":"abc123"},"extra"]`),
			wantReq:   "SUBID",
			wantEvent: []byte(`{"id":"abc123"}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotReq, gotEvent, err := FindEventWithReq(tc.env)

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
			assert.Equal(t, tc.wantReq, gotReq)
			assert.Equal(t, tc.wantEvent, gotEvent)
		})
	}
}

func TestFindSubscriptionEvent(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantSubID   string
		wantEvent   []byte
		wantErr     error
		wantErrText string
	}{
		{
			name:      "valid event",
			env:       []byte(`["EVENT","sub1",{"id":"abc123","kind":1}]`),
			wantSubID: "sub1",
			wantEvent: []byte(`{"id":"abc123","kind":1}`),
		},
		{
			name:        "wrong label",
			env:         []byte(`["REQ","sub1",{"id":"abc123","kind":1}]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected EVENT, got REQ",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["EVENT","sub1"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 3 elements, got 2",
		},
		{
			name:      "extraneous elements",
			env:       []byte(`["EVENT","sub1",{"id":"abc123"},"extra"]`),
			wantSubID: "sub1",
			wantEvent: []byte(`{"id":"abc123"}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotSubID, gotEvent, err := FindSubscriptionEvent(tc.env)

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
			assert.Equal(t, tc.wantSubID, gotSubID)
			assert.Equal(t, tc.wantEvent, gotEvent)
		})
	}
}

func TestFindOK(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantEventID string
		wantStatus  bool
		wantMessage string
		wantErr     error
		wantErrText string
	}{
		{
			name:        "accepted event",
			env:         []byte(`["OK","abc123",true,"Event accepted"]`),
			wantEventID: "abc123",
			wantStatus:  true,
			wantMessage: "Event accepted",
		},
		{
			name:        "rejected event",
			env:         []byte(`["OK","xyz789",false,"Invalid signature"]`),
			wantEventID: "xyz789",
			wantStatus:  false,
			wantMessage: "Invalid signature",
		},
		{
			name:        "wrong status type",
			env:         []byte(`["OK","abc123","ok","Event accepted"]`),
			wantErr:     errors.WrongFieldType,
			wantErrText: "status is not the expected type",
		},
		{
			name:        "wrong label",
			env:         []byte(`["EVENT","abc123",true,"Event accepted"]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected OK, got EVENT",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["OK","abc123",true]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 4 elements, got 3",
		},
		{
			name:        "extraneous elements",
			env:         []byte(`["OK","abc123",true,"Event accepted","extra"]`),
			wantEventID: "abc123",
			wantStatus:  true,
			wantMessage: "Event accepted",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotEventID, gotStatus, gotMessage, err := FindOK(tc.env)

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
			assert.Equal(t, tc.wantEventID, gotEventID)
			assert.Equal(t, tc.wantStatus, gotStatus)
			assert.Equal(t, tc.wantMessage, gotMessage)
		})
	}
}

func TestFindReq(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantSubID   string
		wantFilters [][]byte
		wantErr     error
		wantErrText string
	}{
		{
			name:        "single filter",
			env:         []byte(`["REQ","sub1",{"kinds":[1],"limit":10}]`),
			wantSubID:   "sub1",
			wantFilters: [][]byte{[]byte(`{"kinds":[1],"limit":10}`)},
		},
		{
			name:      "multiple filters",
			env:       []byte(`["REQ","sub2",{"kinds":[1]},{"authors":["abc"]}]`),
			wantSubID: "sub2",
			wantFilters: [][]byte{
				[]byte(`{"kinds":[1]}`),
				[]byte(`{"authors":["abc"]}`),
			},
		},
		{
			name:        "no filters",
			env:         []byte(`["REQ","sub3"]`),
			wantSubID:   "sub3",
			wantFilters: [][]byte{},
		},
		{
			name:        "wrong label",
			env:         []byte(`["EVENT","sub1",{"kinds":[1],"limit":10}]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected REQ, got EVENT",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["REQ"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 2 elements, got 1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotSubID, gotFilters, err := FindReq(tc.env)

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
			assert.Equal(t, tc.wantSubID, gotSubID)
			assert.Equal(t, tc.wantFilters, gotFilters)
		})
	}
}

func TestFindEOSE(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantSubID   string
		wantErr     error
		wantErrText string
	}{
		{
			name:      "valid EOSE",
			env:       []byte(`["EOSE","sub1"]`),
			wantSubID: "sub1",
		},
		{
			name:        "wrong label",
			env:         []byte(`["EVENT","sub1"]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected EOSE, got EVENT",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["EOSE"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 2 elements, got 1",
		},
		{
			name:      "extraneous elements",
			env:       []byte(`["EOSE","sub1","extra"]`),
			wantSubID: "sub1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotSubID, err := FindEOSE(tc.env)

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
			assert.Equal(t, tc.wantSubID, gotSubID)
		})
	}
}

func TestFindClose(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantSubID   string
		wantErr     error
		wantErrText string
	}{
		{
			name:      "valid CLOSE",
			env:       []byte(`["CLOSE","sub1"]`),
			wantSubID: "sub1",
		},
		{
			name:        "wrong label",
			env:         []byte(`["EVENT","sub1"]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected CLOSE, got EVENT",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["CLOSE"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 2 elements, got 1",
		},
		{
			name:      "extraneous elements",
			env:       []byte(`["CLOSE","sub1","extra"]`),
			wantSubID: "sub1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotSubID, err := FindClose(tc.env)

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
			assert.Equal(t, tc.wantSubID, gotSubID)
		})
	}
}

func TestFindClosed(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantSubID   string
		wantMessage string
		wantErr     error
		wantErrText string
	}{
		{
			name:        "valid CLOSED",
			env:         []byte(`["CLOSED","sub1","Subscription complete"]`),
			wantSubID:   "sub1",
			wantMessage: "Subscription complete",
		},
		{
			name:        "wrong label",
			env:         []byte(`["EVENT","sub1","Subscription complete"]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected CLOSED, got EVENT",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["CLOSED","sub1"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 3 elements, got 2",
		},
		{
			name:        "extraneous elements",
			env:         []byte(`["CLOSED","sub1","Subscription complete","extra"]`),
			wantSubID:   "sub1",
			wantMessage: "Subscription complete",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotSubID, gotMessage, err := FindClosed(tc.env)

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
			assert.Equal(t, tc.wantSubID, gotSubID)
			assert.Equal(t, tc.wantMessage, gotMessage)
		})
	}
}

func TestFindNotice(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantMessage string
		wantErr     error
		wantErrText string
	}{
		{
			name:        "valid NOTICE",
			env:         []byte(`["NOTICE","This is a notice"]`),
			wantMessage: "This is a notice",
		},
		{
			name:        "wrong label",
			env:         []byte(`["EVENT","This is a notice"]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected NOTICE, got EVENT",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["NOTICE"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 2 elements, got 1",
		},
		{
			name:        "extraneous elements",
			env:         []byte(`["NOTICE","This is a notice","extra"]`),
			wantMessage: "This is a notice",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotMessage, err := FindNotice(tc.env)

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
			assert.Equal(t, tc.wantMessage, gotMessage)
		})
	}
}

func TestFindAuthChallenge(t *testing.T) {
	cases := []struct {
		name          string
		env           Envelope
		wantChallenge string
		wantErr       error
		wantErrText   string
	}{
		{
			name:          "valid AUTH challenge",
			env:           []byte(`["AUTH","random-challenge-string"]`),
			wantChallenge: "random-challenge-string",
		},
		{
			name:        "wrong label",
			env:         []byte(`["EVENT","random-challenge-string"]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected AUTH, got EVENT",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["AUTH"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 2 elements, got 1",
		},
		{
			name:          "extraneous elements",
			env:           []byte(`["AUTH","random-challenge-string","extra"]`),
			wantChallenge: "random-challenge-string",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotChallenge, err := FindAuthChallenge(tc.env)

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
			assert.Equal(t, tc.wantChallenge, gotChallenge)
		})
	}
}

func TestFindAuthResponse(t *testing.T) {
	cases := []struct {
		name        string
		env         Envelope
		wantEvent   []byte
		wantErr     error
		wantErrText string
	}{
		{
			name:      "valid AUTH response",
			env:       []byte(`["AUTH",{"id":"abc123","kind":22242}]`),
			wantEvent: []byte(`{"id":"abc123","kind":22242}`),
		},
		{
			name:        "wrong label",
			env:         []byte(`["EVENT",{"id":"abc123","kind":22242}]`),
			wantErr:     errors.WrongEnvelopeLabel,
			wantErrText: "expected AUTH, got EVENT",
		},
		{
			name:    "invalid json",
			env:     []byte(`invalid`),
			wantErr: errors.InvalidJSON,
		},
		{
			name:        "missing elements",
			env:         []byte(`["AUTH"]`),
			wantErr:     errors.InvalidEnvelope,
			wantErrText: "expected 2 elements, got 1",
		},
		{
			name:      "extraneous elements",
			env:       []byte(`["AUTH",{"id":"abc123","kind":22242},"extra"]`),
			wantEvent: []byte(`{"id":"abc123","kind":22242}`),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotEvent, err := FindAuthResponse(tc.env)

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
			assert.Equal(t, tc.wantEvent, gotEvent)
		})
	}
}
