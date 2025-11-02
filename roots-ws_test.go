package roots_ws

import "testing"

func TestConnectionStatusConstants(t *testing.T) {
	seen := make(map[ConnectionStatus]bool)

	constants := []ConnectionStatus{
		StatusDisconnected,
		StatusConnecting,
		StatusConnected,
		StatusClosing,
	}

	for i, status := range constants {
		if seen[status] {
			t.Errorf("Duplicate value found for constant at index %d: %d", i, status)
		}
		seen[status] = true
	}
}
