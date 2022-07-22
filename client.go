package gosc

import (
	"fmt"
	"net"
	"regexp"
	"sync"
)

// A Client is an OSC client.
type Client struct {
	remote                 net.Addr
	transport              Transport
	messageHandlers        map[string]MessageHandler
	pendingMessageRequests sync.Map
	bundleHandler          BundleHandler
}

// A MessageHandler is called when messages on the specified address is received.
type MessageHandler interface {
	HandleMessage(writer *ResponseWriter, msg *Message)
}

// A BundleHandler is called when a Bundle is received.
type BundleHandler interface {
	HandleBundle(writer *ResponseWriter, bundle *Bundle)
}

// The MessageHandlerFunc type is an adapter to allow the use of ordinary
// functions as MessageHandler:s. If f is a function with the appropriate
// signature, MessageHandlerFunc(f) is a MessageHandler that calls f.
type MessageHandlerFunc func(writer *ResponseWriter, msg *Message)

// NewClient returns a default client with UDP transport to the given address.
// The client will also start a go-routine to listen for data responses.
//
// The address must be a valid UDP-address including port number.
func NewClient(address string) (*Client, error) {
	remote, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	trans, err := NewUDPTransport(address, 512)
	if err != nil {
		return nil, err
	}

	cli := &Client{
		remote:          remote,
		transport:       trans,
		messageHandlers: map[string]MessageHandler{},
	}
	go cli.listen()

	return cli, nil
}

// HandleMessage adds a MessageHandler for messages on a specific address using
// regexp matching.
func (c *Client) HandleMessage(addressPattern string, handler MessageHandler) error {
	_, err := regexp.Compile(addressPattern)
	if err != nil {
		return fmt.Errorf("addressPattern is not a valid regexp string: %v", err)
	}
	c.messageHandlers[addressPattern] = handler
	return nil
}

// HandleMessageFunc adds a MessageHandlerFunc for messages on a specific
// address using regexp matching.
func (c *Client) HandleMessageFunc(addressPattern string, handlerFunc MessageHandlerFunc) error {
	return c.HandleMessage(addressPattern, handlerFunc)
}

// HandleMessage calls f(msg)
func (m MessageHandlerFunc) HandleMessage(w *ResponseWriter, msg *Message) {
	m(w, msg)
}

// SendMessage uses the clients transport to encode and send an OSC Message
func (c *Client) SendMessage(msg *Message) error {
	return c.transport.Send(msg, c.remote)
}

// SendBundle uses the clients transport to encode and send an OSC Bundle
func (c *Client) SendBundle(bun *Bundle) error {
	return c.transport.Send(bun, c.remote)
}

// SendAndReceiveMessage sends the OSC Message using the clients transport and
// then waits (blocking) for the response to arrive on the listener.
func (c *Client) SendAndReceiveMessage(msg *Message) (*Message, error) {
	ch := make(chan *Message)
	c.pendingMessageRequests.Store(msg.Address, ch)
	err := c.SendMessage(msg)
	if err != nil {
		return nil, err
	}

	res := <-ch
	return res, err
}

// EmitMessage creates an OSC Message using the provided data and then sends it.
func (c *Client) EmitMessage(address string, varArg ...any) error {
	return c.SendMessage(&Message{
		Address:   address,
		Arguments: varArg,
	})
}

// CallMessage creates an OSC Message using the provided data and then sends
// it. See SendAndReceiveMessage
func (c *Client) CallMessage(address string, varArg ...any) (*Message, error) {
	return c.SendAndReceiveMessage(&Message{
		Address:   address,
		Arguments: varArg,
	})
}

func (c *Client) listen() {
	var pkg Package
	var err error
	for pkg, _, err = c.transport.Receive(); err == nil; pkg, _, err = c.transport.Receive() {
		if pkg.GetType() == PackageTypeMessage {
			m := pkg.(*Message)
			if chi, ok := c.pendingMessageRequests.LoadAndDelete(m.Address); ok {
				ch := chi.(chan *Message)
				ch <- m
				close(ch)
			} else {
				for pattern, h := range c.messageHandlers {
					if c.addressMatches(pattern, m.Address) {
						h.HandleMessage(nil, m)
						break
					}
				}
			}
		} else if pkg.GetType() == PackageTypeBundle {
			if c.bundleHandler != nil {
				b := pkg.(*Bundle)
				c.bundleHandler.HandleBundle(nil, b)
			}
		}
	}
}

func (c *Client) addressMatches(pattern, address string) bool {
	matches, _ := regexp.MatchString(pattern, address)
	return matches
}
