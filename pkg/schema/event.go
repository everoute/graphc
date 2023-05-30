package schema

import "encoding/json"

// MutationEvent is the event subscribed from tower
type MutationEvent struct {
	Mutation       MutationType    `json:"mutation"`
	PreviousValues json.RawMessage `json:"previousValues"`
	Node           json.RawMessage `json:"node"`
}

type MutationType string

const (
	CreateEvent MutationType = "CREATED"
	DeleteEvent MutationType = "DELETED"
	UpdateEvent MutationType = "UPDATED"
)
