package client

const (
	evtParameterList string = "AppData,Application,AccountCode,CallerIDName,CallerIDNum,ConnectedLineName,ConnectedLineNum,Context,Exten,Language,Linkedid,Priority,SystemName"
)

// Options describes the options for connecting to
// a native Asterisk AMI- server.
type Options struct {
	//Chatty
	Verbose *bool

	// Usually localhost
	Host string

	//Port to open the asterisk server
	Port int

	// Username for ARI authentication
	Username string

	// Password for ARI authentication
	Password string

	// Allow subscribe to all events in Asterisk Server
	Events string
}

type AmiClient struct {
}
