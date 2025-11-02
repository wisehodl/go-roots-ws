package envelope

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncloseEvent(t *testing.T) {
	cases := []struct {
		name  string
		event []byte
		want  Envelope
	}{
		{
			name:  "empty event",
			event: []byte("{}"),
			want:  []byte(`["EVENT",{}]`),
		},
		{
			name:  "invalid json",
			event: []byte("in[valid,]"),
			want:  []byte(`["EVENT",in[valid,]]`),
		},
		{
			name:  "populated event",
			event: []byte(`{"id":"abc123","kind":1,"sig":"abc123"}`),
			want:  []byte(`["EVENT",{"id":"abc123","kind":1,"sig":"abc123"}]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseEvent(tc.event)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncloseOK(t *testing.T) {
	cases := []struct {
		name    string
		eventID string
		status  bool
		message string
		want    Envelope
	}{
		{
			name:    "successful event",
			eventID: "abc123",
			status:  true,
			message: "Event accepted",
			want:    []byte(`["OK","abc123",true,"Event accepted"]`),
		},
		{
			name:    "rejected event",
			eventID: "xyz789",
			status:  false,
			message: "Invalid signature",
			want:    []byte(`["OK","xyz789",false,"Invalid signature"]`),
		},
		{
			name:    "empty message",
			eventID: "def456",
			status:  true,
			message: "",
			want:    []byte(`["OK","def456",true,""]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseOK(tc.eventID, tc.status, tc.message)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncloseReq(t *testing.T) {
	cases := []struct {
		name    string
		subID   string
		filters [][]byte
		want    Envelope
	}{
		{
			name:    "single filter",
			subID:   "sub1",
			filters: [][]byte{[]byte(`{"kinds":[1],"limit":10}`)},
			want:    []byte(`["REQ","sub1",{"kinds":[1],"limit":10}]`),
		},
		{
			name:    "multiple filters",
			subID:   "sub2",
			filters: [][]byte{[]byte(`{"kinds":[1]}`), []byte(`{"authors":["abc"]}`)},
			want:    []byte(`["REQ","sub2",{"kinds":[1]},{"authors":["abc"]}]`),
		},
		{
			name:    "no filters",
			subID:   "sub3",
			filters: [][]byte{},
			want:    []byte(`["REQ","sub3"]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseReq(tc.subID, tc.filters)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncloseSubscriptionEvent(t *testing.T) {
	cases := []struct {
		name  string
		subID string
		event []byte
		want  Envelope
	}{
		{
			name:  "basic event",
			subID: "sub1",
			event: []byte(`{"id":"abc123","kind":1}`),
			want:  []byte(`["EVENT","sub1",{"id":"abc123","kind":1}]`),
		},
		{
			name:  "empty event",
			subID: "sub2",
			event: []byte(`{}`),
			want:  []byte(`["EVENT","sub2",{}]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseSubscriptionEvent(tc.subID, tc.event)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncloseEOSE(t *testing.T) {
	cases := []struct {
		name  string
		subID string
		want  Envelope
	}{
		{
			name:  "valid subscription ID",
			subID: "sub1",
			want:  []byte(`["EOSE","sub1"]`),
		},
		{
			name:  "empty subscription ID",
			subID: "",
			want:  []byte(`["EOSE",""]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseEOSE(tc.subID)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncloseClose(t *testing.T) {
	cases := []struct {
		name  string
		subID string
		want  Envelope
	}{
		{
			name:  "valid subscription ID",
			subID: "sub1",
			want:  []byte(`["CLOSE","sub1"]`),
		},
		{
			name:  "empty subscription ID",
			subID: "",
			want:  []byte(`["CLOSE",""]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseClose(tc.subID)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncloseClosed(t *testing.T) {
	cases := []struct {
		name    string
		subID   string
		message string
		want    Envelope
	}{
		{
			name:    "with message",
			subID:   "sub1",
			message: "Subscription complete",
			want:    []byte(`["CLOSED","sub1","Subscription complete"]`),
		},
		{
			name:    "empty message",
			subID:   "sub2",
			message: "",
			want:    []byte(`["CLOSED","sub2",""]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseClosed(tc.subID, tc.message)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncloseNotice(t *testing.T) {
	cases := []struct {
		name    string
		message string
		want    Envelope
	}{
		{
			name:    "valid message",
			message: "This is a notice",
			want:    []byte(`["NOTICE","This is a notice"]`),
		},
		{
			name:    "empty message",
			message: "",
			want:    []byte(`["NOTICE",""]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseNotice(tc.message)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncloseAuthChallenge(t *testing.T) {
	cases := []struct {
		name      string
		challenge string
		want      Envelope
	}{
		{
			name:      "valid challenge",
			challenge: "random-challenge-string",
			want:      []byte(`["AUTH","random-challenge-string"]`),
		},
		{
			name:      "empty challenge",
			challenge: "",
			want:      []byte(`["AUTH",""]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseAuthChallenge(tc.challenge)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncloseAuthResponse(t *testing.T) {
	cases := []struct {
		name  string
		event []byte
		want  Envelope
	}{
		{
			name:  "valid event",
			event: []byte(`{"id":"abc123","kind":22242}`),
			want:  []byte(`["AUTH",{"id":"abc123","kind":22242}]`),
		},
		{
			name:  "empty event",
			event: []byte(`{}`),
			want:  []byte(`["AUTH",{}]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncloseAuthResponse(tc.event)
			assert.Equal(t, tc.want, got)
		})
	}
}
