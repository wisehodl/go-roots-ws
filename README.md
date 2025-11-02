# Go-Roots-WS - Nostr WebSocket Transport for Golang

Source: https://git.wisehodl.dev/jay/go-roots-ws

Mirror: https://github.com/wisehodl/go-roots-ws

## What this library does

`go-roots-ws` is a consensus-layer Nostr protocol websocket transport library for golang. It only provides primitives for working with Nostr protocol websocket connection states and messages:

- Websocket Connection States
- Envelope Structure
- Message Validation
- Protocol Message Creation
- Protocol Message Parsing
- Standard Label Handling

## What this library does not do

`go-roots-ws` serves as a foundation for other libraries and applications to implement higher level transport abstractions on top of it, including:

- Connection Management
- Event Loops
- Subscription Handling
- State Management
- Reconnection Logic

## Installation

1. Add `go-roots-ws` to your project:

```bash
go get git.wisehodl.dev/jay/go-roots-ws
```

If the primary repository is unavailable, use the `replace` directive in your go.mod file to get the package from the github mirror:

```
replace git.wisehodl.dev/jay/go-roots-ws => github.com/wisehodl/go-roots-ws latest
```

2. Import the packages:

```golang
import (
    "encoding/json"
    "git.wisehodl.dev/jay/go-roots/events"
    "git.wisehodl.dev/jay/go-roots/filters"
    "git.wisehodl.dev/jay/go-roots-ws/envelope"
    "git.wisehodl.dev/jay/go-roots-ws/errors"
)
```

3. Access functions with appropriate namespaces.

## Usage Examples

### Envelope Creation

#### Create EVENT envelope

```go
// Create an event using go-roots
event := events.Event{
    ID:        "abc123",
    PubKey:    "def456",
    Kind:      1,
    Content:   "Hello Nostr!",
    CreatedAt: int(time.Now().Unix()),
}

// Convert to JSON
eventJSON, err := json.Marshal(event)
if err != nil {
    log.Fatal(err)
}

// Create envelope
env := envelope.EncloseEvent(eventJSON)
// Result: ["EVENT",{"id":"abc123","pubkey":"def456","kind":1,"content":"Hello Nostr!","created_at":1636394097}]
```

#### Create subscription EVENT envelope

```go
// Create an event using go-roots
event := events.Event{
    ID:        "abc123",
    PubKey:    "def456",
    Kind:      1,
    Content:   "Hello Nostr!",
    CreatedAt: int(time.Now().Unix()),
}

// Convert to JSON
eventJSON, err := json.Marshal(event)
if err != nil {
    log.Fatal(err)
}

// Create envelope with subscription ID
subID := "sub1"
env := envelope.EncloseSubscriptionEvent(subID, eventJSON)
// Result: ["EVENT","sub1",{"id":"abc123","pubkey":"def456","kind":1,"content":"Hello Nostr!","created_at":1636394097}]
```

#### Create REQ envelope

```go
// Create filters using go-roots
since := int(time.Now().Add(-24 * time.Hour).Unix())
limit := 50

filter1 := filters.Filter{
    Kinds: []int{1},
    Limit: &limit,
    Since: &since,
}

filter2 := filters.Filter{
    Authors: []string{"def456"},
}

// Marshal filters to JSON
filter1JSON, err := filters.MarshalJSON(filter1)
if err != nil {
    log.Fatal(err)
}

filter2JSON, err := filters.MarshalJSON(filter2)
if err != nil {
    log.Fatal(err)
}

// Create envelope
subID := "sub1"
filtersJSON := [][]byte{filter1JSON, filter2JSON}
env := envelope.EncloseReq(subID, filtersJSON)
// Result: ["REQ","sub1",{"kinds":[1],"limit":50,"since":1636307697},{"authors":["def456"]}]
```

#### Create other envelope types

```go
// Create CLOSE envelope
env := envelope.EncloseClose("sub1")
// Result: ["CLOSE","sub1"]

// Create EOSE envelope
env := envelope.EncloseEOSE("sub1")
// Result: ["EOSE","sub1"]

// Create NOTICE envelope
env := envelope.EncloseNotice("This is a notice")
// Result: ["NOTICE","This is a notice"]

// Create OK envelope
env := envelope.EncloseOK("abc123", true, "Event accepted")
// Result: ["OK","abc123",true,"Event accepted"]

// Create AUTH challenge
env := envelope.EncloseAuthChallenge("random-challenge-string")
// Result: ["AUTH","random-challenge-string"]

// Create AUTH response
// Create an event using go-roots
authEvent := events.Event{
    ID:        "abc123",
    PubKey:    "def456",
    Kind:      22242,
    Content:   "",
    CreatedAt: int(time.Now().Unix()),
}

// Convert to JSON
authEventJSON, err := json.Marshal(authEvent)
if err != nil {
    log.Fatal(err)
}

// Create envelope
env := envelope.EncloseAuthResponse(authEventJSON)
// Result: ["AUTH",{"id":"abc123","pubkey":"def456","kind":22242,"content":"","created_at":1636394097}]
```

---

### Envelope Parsing

#### Extract label from envelope

```go
env := []byte(`["EVENT",{"id":"abc123","pubkey":"def456","kind":1,"content":"Hello Nostr!"}]`)
label, err := envelope.GetLabel(env)
if err != nil {
    log.Fatal(err)
}
// label: "EVENT"

// Check if label is standard
isStandard := envelope.IsStandardLabel(label)
// isStandard: true
```

#### Extract event from EVENT envelope

```go
env := []byte(`["EVENT",{"id":"abc123","pubkey":"def456","kind":1,"content":"Hello Nostr!"}]`)
eventJSON, err := envelope.FindEvent(env)
if err != nil {
    log.Fatal(err)
}

// Parse into go-roots Event
var event events.Event
err = json.Unmarshal(eventJSON, &event)
if err != nil {
    log.Fatal(err)
}

// Validate the event
if err := events.Validate(event); err != nil {
    log.Printf("Invalid event: %v", err)
}

// Now you can access event properties
fmt.Println(event.ID, event.Kind, event.Content)
```

#### Extract subscription event

```go
env := []byte(`["EVENT","sub1",{"id":"abc123","pubkey":"def456","kind":1,"content":"Hello Nostr!"}]`)
subID, eventJSON, err := envelope.FindSubscriptionEvent(env)
if err != nil {
    log.Fatal(err)
}

// Parse into go-roots Event
var event events.Event
err = json.Unmarshal(eventJSON, &event)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Subscription: %s, Event ID: %s\n", subID, event.ID)
```

#### Extract subscription request

```go
env := []byte(`["REQ","sub1",{"kinds":[1],"limit":50},{"authors":["def456"]}]`)
subID, filtersJSON, err := envelope.FindReq(env)
if err != nil {
    log.Fatal(err)
}

// Parse each filter
var parsedFilters []filters.Filter
for _, filterJSON := range filtersJSON {
    var filter filters.Filter
    err := filters.UnmarshalJSON(filterJSON, &filter)
    if err != nil {
        log.Fatal(err)
    }
    parsedFilters = append(parsedFilters, filter)
}

// Now you can use the filter objects
for i, filter := range parsedFilters {
    fmt.Printf("Filter %d: %+v\n", i, filter)
}
```

#### Extract other envelope types

```go
// Extract OK response
env := []byte(`["OK","abc123",true,"Event accepted"]`)
eventID, status, message, err := envelope.FindOK(env)
// eventID: "abc123"
// status: true
// message: "Event accepted"

// Extract EOSE message
env := []byte(`["EOSE","sub1"]`)
subID, err := envelope.FindEOSE(env)
// subID: "sub1"

// Extract CLOSE message
env := []byte(`["CLOSE","sub1"]`)
subID, err := envelope.FindClose(env)
// subID: "sub1"

// Extract CLOSED message
env := []byte(`["CLOSED","sub1","Subscription complete"]`)
subID, message, err := envelope.FindClosed(env)
// subID: "sub1"
// message: "Subscription complete"

// Extract NOTICE message
env := []byte(`["NOTICE","This is a notice"]`)
message, err := envelope.FindNotice(env)
// message: "This is a notice"

// Extract AUTH challenge
env := []byte(`["AUTH","random-challenge-string"]`)
challenge, err := envelope.FindAuthChallenge(env)
// challenge: "random-challenge-string"

// Extract AUTH response
env := []byte(`["AUTH",{"id":"abc123","pubkey":"def456","kind":22242,"content":""}]`)
authEventJSON, err := envelope.FindAuthResponse(env)
if err != nil {
    log.Fatal(err)
}

// Parse into go-roots Event
var authEvent events.Event
err = json.Unmarshal(authEventJSON, &authEvent)
if err != nil {
    log.Fatal(err)
}
```

## Testing

This library contains a comprehensive suite of unit tests. Run them with:

```bash
go test ./...
```
