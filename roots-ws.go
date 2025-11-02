package roots_ws

// ConnectionStatus represents the current state of a WebSocket connection.
type ConnectionStatus int

const (
	// StatusDisconnected indicates the connection is not active and no connection attempt is in progress.
	StatusDisconnected ConnectionStatus = iota

	// StatusConnecting indicates a connection attempt is currently in progress but not yet established.
	StatusConnecting

	// StatusConnected indicates the connection is active and ready for message exchange.
	StatusConnected

	// StatusClosing indicates the connection is in the process of shutting down gracefully.
	StatusClosing
)
