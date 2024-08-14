package responses

import "github.com/cripito/amilib"

type Response interface {
	// GetPayload returns the content of this Response
	GetPayload() string

	// GetNode returns the unique ID of the Asterisk system on which this Response originated
	GetNode() string

	// GetType returns the type name of this Response
	GetType() ResponseType

	//GetActionID returns the uniqueID for the response
	GetID() string

	//GetMessage returns the messages  for the response
	GetMessage() string
}

type ResponseType int

const (
	Success = iota
	Error
)

func (w ResponseType) String() string {
	return [...]string{"Success", "Errors"}[w]
}

type ResponseData struct {
	// ActionId of the response
	ID string `json:"actionid,omitempty"`

	// Payload indicates the content of this Response
	PayLoad string `json:"payload,omitempty"`

	// Node indicates the unique identifier of the source Asterisk box for this Response
	Node string `json:"asterisk_id,omitempty"`

	// Timestamp indicates the time this Response was generated
	Timestamp amilib.DateTime `json:"timestamp,omitempty"`

	// Type is the type name of this response //Success-Error
	Type ResponseType `json:"type"`

	// Message of the response
	Message string `json:"message,omitempty"`
}

func NewResponseData() *ResponseData {
	return &ResponseData{}
}

// GetType gets the type of the event
func (e *ResponseData) GetType() ResponseType {
	return e.Type
}

// GetNode gets the node ID of the source Asterisk instance
func (e *ResponseData) GetNode() string {
	return e.Node
}

func (e *ResponseData) GetPayload() string {
	return e.PayLoad
}

func (e *ResponseData) GetID() string {
	return e.ID
}

func (e *ResponseData) GetMessage() string {
	return e.Message
}
