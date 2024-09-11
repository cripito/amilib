package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/cripito/amilib/logs"
	"github.com/cripito/amilib/responses"
	"github.com/cripito/amilib/rid"
	amitools "github.com/cripito/amilib/tools"
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

	//Request
	requests map[string]*amitools.Request

	// List of natsubjects the system will process
	Node *amitools.Node

	//Nats subjects
	subsSujects []string

	//nats subscriptions
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

	mu sync.Mutex

	nodelist map[string]*amitools.Node

	amiRequestHandler AMIRequestHandler
}

var (
	registeredHandlers sync.Map
)

func init() {
	registeredHandlers = sync.Map{}
}

type key[TRequest any, TResult any] struct{}

type RequestHandler[TRequest any, TResult any] func(req TRequest) (TResult, error)

type AMIRequestHandler func(req *amitools.Request) (*amitools.Request, error)

// OriginateData holds information used to originate outgoing calls.
//
//	Channel - Channel name to call.
//	Exten - Extension to use (requires Context and Priority)
//	Context - Context to use (requires Exten and Priority)
//	Priority - Priority to use (requires Exten and Context)
//	Application - Application to execute.
//	Data - Data to use (requires Application).
//	Timeout - How long to wait for call to be answered (in ms.).
//	CallerID - Caller ID to be set on the outgoing channel.
//	Variable - Channel variable to set, multiple Variable: headers are allowed.
//	Account - Account code.
//	EarlyMedia - Set to true to force call bridge on early media.
//	Async - Set to true for fast origination.
//	Codecs - Comma-separated list of codecs to use for this call.
//	ChannelId - Channel UniqueId to be set on the channel.
//	OtherChannelId - Channel UniqueId to be set on the second local channel.
type OriginateData struct {
	Channel        string   `ami:"Channel,omitempty"`
	Exten          string   `ami:"Exten,omitempty"`
	Context        string   `ami:"Context,omitempty"`
	Priority       int      `ami:"Priority,omitempty"`
	Application    string   `ami:"Application,omitempty"`
	Data           string   `ami:"Data,omitempty"`
	Timeout        int      `ami:"Timeout,omitempty"`
	CallerID       string   `ami:"CallerID,omitempty"`
	Variable       []string `ami:"Variable,omitempty"`
	Account        string   `ami:"Account,omitempty"`
	EarlyMedia     string   `ami:"EarlyMedia,omitempty"`
	Async          string   `ami:"Async,omitempty"`
	Codecs         string   `ami:"Codecs,omitempty"`
	ChannelID      string   `ami:"ChannelId,omitempty"`
	OtherChannelID string   `ami:"OtherChannelId,omitempty"`
}

func NewAmiClient(ctx context.Context, opts ...OptionFunc) *AMIClient {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logs.Init()

	ami := &AMIClient{
		MBPrefix:       "ami.",
		TopicSeparator: ".",
		MBPostfix:      ".*",
		serverName:     "ami.client.",
		Node: &amitools.Node{
			Type: amitools.NodeType_UNSPECIFIED,
		},
		requests: make(map[string]*amitools.Request),
		nodelist: make(map[string]*amitools.Node),
	}

	// Load explicit configurations
	for _, opt := range opts {
		opt(ami)
	}

	if ami.serverName == "" {
		ami.serverName = "prx"
	}

	ami.setNode()

	return ami
}

func (s *AMIClient) SetAmiRequestHandler(fcall AMIRequestHandler) {
	s.amiRequestHandler = fcall
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

func WithType(t amitools.NodeType) OptionFunc {
	return func(c *AMIClient) {
		c.Node.Type = t
	}
}

func (ami *AMIClient) GetNats() *nats.Conn {
	return ami.mbus
}

func (ami *AMIClient) GetNode() *amitools.Node {
	return ami.Node
}

func (ami *AMIClient) setNode() error {
	var err error

	prefix := ""
	if ami.Node.Type == amitools.NodeType_PROXY {
		prefix = "prx"
	} else {
		prefix = "cli"
	}

	ami.Node.ID, err = rid.New(prefix)
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return err
	}

	if ami.Node.Type == amitools.NodeType_PROXY {
		ami.Node.Name = ami.MBPrefix + "proxy" + ami.TopicSeparator + ami.Node.ID
		ami.Node.IP = "ip"

	} else {
		ami.Node.Name = ami.MBPrefix + "client" + ami.TopicSeparator + ami.Node.ID
		ami.Node.Type = amitools.NodeType_CLIENT
	}

	ami.Node.Status = amitools.StatusType_UNKNOWN

	ami.nodelist[ami.Node.ID] = ami.Node

	return nil
}

func (ami *AMIClient) getNodeByID(id string) *amitools.Node {
	if node, found := ami.nodelist[id]; found {
		return node
	} else {
		return nil
	}
}

func (ami *AMIClient) GetNodeByIP(ip string) *amitools.Node {
	for _, node := range ami.nodelist {
		if node.IP == ip {
			return node
		}
	}

	return nil
}

func (ami *AMIClient) natsConnection(ctx context.Context) error {
	var err error

	logs.TLogger.Debug().Msgf("Connecting to NATS using %s with name %s", ami.mbusURI, ami.Node.ID)

	ami.mbus, err = nats.Connect(ami.mbusURI,
		nats.Name(ami.Node.Name),
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

func (ami *AMIClient) announceHandler(msg *nats.Msg) {
	logs.TLogger.Debug().Msgf("GOT announce %s", msg.Data)
	node := &amitools.Node{}

	err := json.Unmarshal(msg.Data, node)
	if err != nil {
		logs.TLogger.Error().Msgf(err.Error())

		return
	}

	inode := ami.getNodeByID(node.ID)
	if inode != nil {
		inode.Status = node.Status
		inode.LastPing = node.LastPing
	} else {
		ami.nodelist[node.ID] = node
	}

}

func (ami *AMIClient) Subjects(topic string, id string) string {
	return fmt.Sprintf("%s%s%s%s", ami.MBPrefix, topic, ami.TopicSeparator, id)
}

func (ami *AMIClient) GetNodes() map[string]*amitools.Node {
	return ami.nodelist
}

func remove(slice []*nats.Subscription, s int) []*nats.Subscription {
	return append(slice[:s], slice[s+1:]...)
}

func (ami *AMIClient) UnSubscribe(sub *nats.Subscription) error {
	if sub != nil {
		for i, tsub := range ami.subs {
			if tsub == sub {
				ami.subs = remove(ami.subs, i)
			}
		}

		return sub.Unsubscribe()
	}

	return errors.New("nil subscription")
}

func (ami *AMIClient) Subscribe(evt *amitools.Event, node *amitools.Node, cb nats.MsgHandler) (*nats.Subscription, error) {

	var sub *nats.Subscription
	var err error

	var subject string = ami.MBPrefix + "events" + ami.TopicSeparator + "*"
	if evt != nil {
		subject = ami.MBPrefix + "events" + ami.TopicSeparator + evt.Type
	}

	if node != nil {
		subject = subject + ami.TopicSeparator + ami.Node.ID
	} else {
		subject = subject + ami.TopicSeparator + "*"
	}

	logs.TLogger.Debug().Msgf("subscribing to %s", subject)
	sub, err = ami.mbus.Subscribe(subject, cb)
	if err != nil {
		return nil, err
	}

	ami.subs = append(ami.subs, sub)

	return sub, nil
}

// register handlers for invocation in the mediator
func Register[TRequest any, TResult any](fcall RequestHandler[TRequest, TResult]) error {
	k := key[TRequest, TResult]{}

	_, existed := registeredHandlers.LoadOrStore(reflect.TypeOf(k), fcall)
	if existed {
		return errors.New("the provided type is already registered to a handler")
	}

	return nil
}

// invoke a register handler
func Invoke[TRequest any, TResult any](req TRequest) (TResult, error) {
	var zeroRes TResult
	var k key[TRequest, TResult]

	handler, ok := registeredHandlers.Load(reflect.TypeOf(k))
	if !ok {
		return zeroRes, errors.New("could not find zeroRes handler for this function")
	}

	switch handler := handler.(type) {
	case RequestHandler[TRequest, TResult]:
		return handler(req)
	}

	return zeroRes, errors.New("Invalid handler")
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

	switch ami.Node.Type {
	case amitools.NodeType_PROXY:
	case amitools.NodeType_CLIENT:
		logs.TLogger.Debug().Msgf("subscribing to %s", ami.Subjects("announce", "*"))
		Sub, err := ami.mbus.Subscribe(ami.Subjects("announce", "*"), ami.announceHandler)
		if err != nil {
			logs.TLogger.Error().Msg(err.Error())

			return err
		}

		ami.subs = append(ami.subs, Sub)
	case amitools.NodeType_UNSPECIFIED:
		err := errors.New("client type not specified")
		if err != nil {
			logs.TLogger.Error().Msg(err.Error())

			return err
		}
	}

	logs.TLogger.Debug().Msgf("subscribing to %s", ami.Subjects("requests", ami.Node.ID))
	tSub, err := ami.mbus.Subscribe(ami.Subjects("requests", ami.Node.ID), func(msg *nats.Msg) {
		logs.TLogger.Debug().Msgf("WE GOT %s from %s", msg.Data, msg.Reply)

		req := &amitools.Request{}

		err := json.Unmarshal(msg.Data, req)
		if err != nil {
			logs.TLogger.Error().Msgf(err.Error())

			return
		}

		switch ami.Node.GetType() {
		case amitools.NodeType_CLIENT:
		case amitools.NodeType_PROXY:
			ami.requests[req.ID] = req

			if ami.amiRequestHandler != nil {
				req, err = ami.amiRequestHandler(req)
				if err != nil {
					logs.TLogger.Error().Msgf(err.Error())

					return
				}
			}
		case amitools.NodeType_UNSPECIFIED:
		}

		data, err := json.Marshal(req)
		if err != nil {
			logs.TLogger.Error().Msg(err.Error())

			return
		}

		ami.mbus.Publish(msg.Reply, data)
	})
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return err
	}
	ami.subs = append(ami.subs, tSub)

	go func() {
		for {
			select {
			case <-time.After(55 * time.Second):
				for r, node := range ami.nodelist {

					ami.mu.Lock()

					switch node.Status {
					case amitools.StatusType_DOWN:
						delete(ami.nodelist, r)
					case amitools.StatusType_UP:
						node.Status = amitools.StatusType_DOWN

						ami.nodelist[r] = node
					case amitools.StatusType_UNKNOWN:
						node.Status = amitools.StatusType_UNKNOWN

						ami.nodelist[r] = node
					}

					ami.mu.Unlock()
				}
			case <-ami.shutdown:
				return
			}
		}

	}()

	return nil
}

func (ami *AMIClient) Close() {

	for _, sub := range ami.subs {
		sub.Unsubscribe()
		sub.Drain()
	}

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

func (ami *AMIClient) read(ctx context.Context, id string) (*responses.ResponseData, error) {

	return nil, nil
}

func (ami *AMIClient) send(ctx context.Context, req *amitools.Request, node *amitools.Node) *amitools.Request {

	data, err := json.Marshal(req)
	if err != nil {
		logs.TLogger.Error().Msg(err.Error())

		return nil
	}

	natsConn := ami.GetNats()
	if natsConn != nil {
		logs.TLogger.Debug().Msgf("Request to %s ", ami.MBPrefix+"requests"+ami.TopicSeparator+node.ID)

		resp, err := natsConn.Request(ami.MBPrefix+"requests"+ami.TopicSeparator+node.ID, data, 3*time.Second)
		if err != nil {
			logs.TLogger.Error().Msg(err.Error())

			return nil
		}

		//logs.TLogger.Debug().Msgf("Answer %s", resp.Data)

		p := &amitools.Request{}

		err = json.Unmarshal(resp.Data, p)
		if err != nil {
			logs.TLogger.Error().Msg(err.Error())

			return nil
		}

		return p
	}

	return nil
}

func (ami *AMIClient) createRequest(action string, id string, v ...interface{}) (*amitools.Request, error) {

	b, err := command(action, id, v)
	if err != nil {
		return nil, err
	}

	p := &amitools.Request{
		ID:      id,
		Type:    action,
		Node:    ami.Node,
		Payload: string(b),
	}

	return p, nil
}

func (ami *AMIClient) CoreSettings(ctx context.Context, actionID string, node *amitools.Node) (*amitools.Request, error) {
	var settings = struct{}{}

	p, err := ami.createRequest("coresettings", actionID, settings)
	if err != nil {
		return nil, err
	}

	resp := ami.send(ctx, p, node)

	return resp, nil
}

func (ami *AMIClient) Originate(ctx context.Context, actionID string, originate *OriginateData, node *amitools.Node) (*amitools.Request, error) {
	p, err := ami.createRequest("Originate", actionID, originate)
	if err != nil {
		return nil, err
	}

	resp := ami.send(ctx, p, node)

	return resp, nil
}

func (ami *AMIClient) Hangup(ctx context.Context, actionID string, channel string, cause string) (*amitools.Request, error) {
	channelMap := map[string]string{
		"Channel": channel,
		"Cause":   cause,
	}
}
