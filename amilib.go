package amilib

import (
	"context"

	"github.com/cripito/amilib/request"
	"github.com/rs/zerolog"
)

type AMIClientInterface interface {
	// Close shuts down the client
	Close() error

	//Connect to the client
	Connect(ctx context.Context) error

	//Set the logger handler for client
	SetLogHandler(log zerolog.Logger)

	//Set the listen for the client
	Listen(ctx context.Context) error
}

type RequestHandler interface {
	CreateRequest(*request.Request)
}

type RequestHandlerFunc func(*request.Request) error

type EventType string

const (
	All               EventType = "all"
	FullyBooted       EventType = "FullyBooted"
	Hangup            EventType = "Hangup"
	Newstate          EventType = "Newstate"
	DialEnd           EventType = "DialEnd"
	SoftHangupRequest EventType = "SoftHangupRequest"
	HangupRequest     EventType = "HangupRequest"
	MixMonitorStart   EventType = "MixMonitorStart"
	MixMonitorStop    EventType = "MixMonitorStop"
	NewConnectedLine  EventType = "NewConnectedLine"
	DialBegin         EventType = "DialBegin"
	NewCallerid       EventType = "NewCallerid"
	NewAccountCode    EventType = "NewAccountCode"
)

func (evt EventType) String() string {
	return string(evt)
}
