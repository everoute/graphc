package informer

import (
	"github.com/everoute/graphc/pkg/schema"
)

type CrcEventType string

const (
	CrcEventInsert CrcEventType = "INSERT"
	CrcEventUpdate CrcEventType = "UPDATE"
	CrcEventDelete CrcEventType = "DELETE"
)

type CrcEvent struct {
	EventType CrcEventType
	OldObj    schema.Object
	NewObj    schema.Object
}
