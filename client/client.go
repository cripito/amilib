package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cripito/amilib/logs"
	"github.com/cripito/amilib/node"
	"github.com/cripito/amilib/request"
	"github.com/cripito/amilib/responses"
	"github.com/cripito/amilib/rid"
	"github.com/nats-io/nats.go"
)

const (
	evtParameterList string = "AppData,Application,AccountCode,CallerIDName,CallerIDNum,ConnectedLineName,ConnectedLineNum,Context,Exten,Language,Linkedid,Priority,SystemName"
)

// Options describes the options for connecting to
// a native Asterisk AMI- server.
type Options struct {
	//Chatty
	Verbose *bool

	MBus *nats.Conn
}

// OptionFunc is a function which configures options on a Client
type OptionFunc func(*AMIClient)

type AMIClient struct {
	// List of natsubjects the system will process
	Node *node.Node

	subsSujects []string

	subs []*nats.Subscription

	// MBPrefix is the string which should be prepended to all Nats subjects, sending and receiving.
	// It defaults to "ami.".
	MBPrefix string

	//TopicSeparator is the string which we use btw topics to all Nats Subjects, sending and receiving.
	TopicSeparator string

	// MBPrefix is the string which should be add to the end to all Nats subjects, sending and receiving.
	// It defaults to ".*".
	MBPostfix string

	//Nats Uri
	mbusURI string

	//nats connection handler
	mbus *nats.Conn

	//nats server connection name
	serverName string

	//natserver connected
	connectedServer string

	shutdown chan os.Signal

	natsfuncHandlerAnnouncement nats.MsgHandler
	natsfuncHandlerRequest      nats.MsgHandler
	natsfuncHandlerEvents       nats.MsgHandler
}

func NewAmiClient(ctx context.Context, opts ...OptionFunc) *AMIClient {

	var err error

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logs.Init()

	ami := &AMIClient{
		MBPrefix:       "ami.",
		TopicSeparator: ".",
		MBPostfix:      ".*",
		serverName:     "ami.client.",
		Node: &node.Node{
			NodeIP: node.WHOCARES,
			Type:   node.Client,
		},
	}

	// Load explicit configurations
	for _, opt := range opts {
		opt(ami)
	}

	ami.Node.NodeID, err = rid.New(ami.serverName)
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return nil
	}

	return ami
}

func WithNatsUri(uri string) OptionFunc {
	return func(c *AMIClient) {
		c.mbusURI = uri
	}
}

func WithName(name string) OptionFunc {
	return func(c *AMIClient) {
		c.serverName = name
	}
}

func WithPrefix(prefix string) OptionFunc {
	return func(c *AMIClient) {
		c.MBPrefix = prefix
	}
}

func (ami *AMIClient) GetNats() *nats.Conn {
	return ami.mbus
}

func (ami *AMIClient) natsConnection(ctx context.Context) error {
	var err error

	logs.TLogger.Debug().Msgf("Connecting to NATS using %s with name %s", ami.mbusURI, ami.Node.NodeID)

	ami.mbus, err = nats.Connect(ami.mbusURI,
		nats.Name(ami.Node.NodeID),
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			logs.TLogger.Debug().Msgf("Known servers: %v", nc.Servers())
			logs.TLogger.Debug().Msgf("Discovered servers: %v", nc.DiscoveredServers())
		}),
		nats.ReconnectBufSize(1*1024*1024),
		nats.ReconnectWait(10*time.Second),
		nats.MaxReconnects(10),
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(3),
		nats.NoEcho(),
		nats.DisconnectHandler(func(c *nats.Conn) {
			logs.TLogger.Debug().Msgf("Disconnected FROM %s", ami.connectedServer)
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			ami.mbus = c
			ami.connectedServer = ami.mbus.ConnectedAddr()

			logs.TLogger.Debug().Msgf("Reconnected to %s", ami.connectedServer)
		}),
	)
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return nil
	}

	return nil
}

func (ami *AMIClient) Subjects(topic string) string {
	return fmt.Sprintf("%s%s%s", ami.MBPrefix, topic, ami.MBPostfix)

}

func (ami *AMIClient) runAnnouncementFunc(ctx context.Context, fcallback nats.MsgHandler) {
	logs.TLogger.Debug().Msgf("subscribing to %s", ami.Subjects("announce"))
	tSub, err := ami.mbus.Subscribe(ami.Subjects("announce"), fcallback)
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return
	}

	ami.subs = append(ami.subs, tSub)
	ami.subsSujects = append(ami.subsSujects, ami.Subjects("announce"))
}

func (ami *AMIClient) runRequestFunc(ctx context.Context, fcallback nats.MsgHandler) {
	logs.TLogger.Debug().Msgf("subscribing to %s", ami.serverName)

	tSub, err := ami.mbus.Subscribe(ami.serverName, fcallback)
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return
	}

	ami.subs = append(ami.subs, tSub)
	ami.subsSujects = append(ami.subsSujects, ami.Subjects("request"))
}

func (ami *AMIClient) runEventFunc(ctx context.Context, fcallback nats.MsgHandler) {
	logs.TLogger.Debug().Msgf("subscribing to %s", ami.Subjects("events"))

	tSub, err := ami.mbus.Subscribe(ami.Subjects("events"), fcallback)
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return
	}

	ami.subs = append(ami.subs, tSub)
	ami.subsSujects = append(ami.subsSujects, ami.Subjects("events"))
}

func (ami *AMIClient) SetAnnouceHandler(fcallback nats.MsgHandler) {
	ami.natsfuncHandlerAnnouncement = fcallback
}

func (ami *AMIClient) SetRequestHandler(fcallback nats.MsgHandler) {
	ami.natsfuncHandlerRequest = fcallback
}

func (ami *AMIClient) SetEventsHandler(fcallback nats.MsgHandler) {
	ami.natsfuncHandlerEvents = fcallback
}

// ----------------------------------------
func (ami *AMIClient) Listen(ctx context.Context) error {

	ami.shutdown = make(chan os.Signal)
	signal.Notify(ami.shutdown, os.Interrupt, syscall.SIGTERM)

	err := ami.natsConnection(ctx)
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return err
	}

	if ami.natsfuncHandlerAnnouncement != nil {
		ami.runAnnouncementFunc(context.Background(), ami.natsfuncHandlerAnnouncement)
	}

	if ami.natsfuncHandlerRequest != nil {
		ami.runRequestFunc(context.Background(), ami.natsfuncHandlerRequest)
	}

	if ami.natsfuncHandlerEvents != nil {
		ami.runEventFunc(context.Background(), ami.natsfuncHandlerEvents)
	}

	for {
		select {
		case <-ami.shutdown:
			return nil
		}
	}

}

func (ami *AMIClient) Close() {
	if ami.mbus != nil {
		ami.mbus.Drain()
		ami.mbus.Close()
	}
}

func command(action string, id string, v ...interface{}) ([]byte, error) {
	if action == "" {
		return nil, errors.New("invalid Action")
	}
	return marshal(&struct {
		Action string `ami:"Action"`
		ID     string `ami:"ActionID, omitempty"`
		V      []interface{}
	}{Action: action, ID: id, V: v})
}

func (ami *AMIClient) send(ctx context.Context, req *request.Request, nodeid string) *responses.ResponseData {
	return nil
}

func (ami *AMIClient) createRequest(action string, id string, v ...interface{}) (*request.Request, error) {

	b, err := command(action, id, v)
	if err != nil {
		return nil, err
	}

	p := &request.Request{
		ID:      id,
		Type:    action,
		Node:    ami.Node,
		Payload: string(b),
	}

	return p, nil
}

func (ami *AMIClient) CoreSettings(ctx context.Context, actionID string, node string) (*responses.ResponseData, error) {
	var settings = struct{}{}

	p, err := ami.createRequest("coresettings", actionID, settings)
	if err != nil {
		return nil, err
	}

	resp := ami.send(ctx, p, node)
	if resp != nil && resp.GetType() == responses.Error {
		err := fmt.Errorf("error: %s", resp.GetMessage())

		logs.TLogger.Debug().Msg(err.Error())

		return nil, err

	}

	return resp, nil
}
