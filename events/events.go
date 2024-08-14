package events

import (
	"github.com/cripito/amilib"
)

// Event is the top level event interface
type Event interface {
	// GetPayload returns any payload by which this event has been tagged
	GetPayload() string

	// GetNode returns the unique ID of the Asterisk system on which this event originated
	GetNode() string

	// GetType returns the type name of this event
	GetType() EventType

	//GetActionID returns the uniqueID for the event
	GetID() string
}

// EventData provides the basic metadata for an AMI event
type EventData struct {
	// ActionId of the response
	ID string `json:"actionid,omitempty"`

	// Payload indicates the content of this event
	PayLoad string `json:"payload,omitempty"`

	// Node indicates the unique identifier of the source Asterisk box for this event
	Node string `json:"asterisk_id,omitempty"`

	// Timestamp indicates the time this event was generated
	Timestamp amilib.DateTime `json:"timestamp,omitempty"`

	// Type is the type name of this event
	Type EventType `json:"type"`
}

func NewEventData() *EventData {
	return &EventData{}
}

// GetType gets the type of the event
func (e *EventData) GetType() EventType {
	return e.Type
}

// GetNode gets the node ID of the source Asterisk instance
func (e *EventData) GetNode() string {
	return e.Node
}

func (e *EventData) GetPayload() string {
	return e.PayLoad
}

func (e *EventData) GetID() string {
	return e.ID
}

type EventType int

const (
	None = iota
	All
)

func (w EventType) String() string {
	return [...]string{"None", "All"}[w]
}
