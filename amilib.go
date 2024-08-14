package amilib

import "context"

type AMIClientInterface interface {
	// Close shuts down the client
	Close()

	//Connect to the client
	Connect(ctx context.Context)
}
