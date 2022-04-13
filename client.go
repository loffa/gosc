package gosc

import (
	"sync"
)

// A Client is an OSC client.
type Client struct {
	transport              Transport
	messageHandlers        map[string]MessageHandler
	pendingMessageRequests sync.Map
	bundleHandler          BundleHandler
}

// A MessageHandler is called when messages on the specified address is received.
type MessageHandler interface {
	HandleMessage(msg *Message)
}

// A BundleHandler is called when a Bundle is received.
type BundleHandler interface {
	HandleBundle(bundle *Bundle)
}

// The MessageHandlerFunc type is an adapter to allow the use of ordinary
// functions as MessageHandler:s. If f is a function with the appropriate
// signature, MessageHandlerFunc(f) is a MessageHandler that calls f.
type MessageHandlerFunc func(msg *Message)

// NewClient returns a default client with UDP transport to the given address.
//
// The address must be a valid UDP-address including port number.
func NewClient(address string) (*Client, error) {
	trans, err := NewUDPTransport(address, 512)
	if err != nil {
		return nil, err
	}
	cli := &Client{
		transport:       trans,
		messageHandlers: map[string]MessageHandler{},
	}
	go cli.listen()

	return cli, nil
}

// HandleMessage adds a MessageHandler for messages on a specific address.
func (c *Client) HandleMessage(address string, handler MessageHandler) {
	c.messageHandlers[address] = handler
}

// HandleMessageFunc adds a MessageHandlerFunc for messages on a specific
// address.
func (c *Client) HandleMessageFunc(address string, handlerFunc MessageHandlerFunc) {
	c.HandleMessage(address, handlerFunc)
}

// HandleMessage calls f(msg)
func (f MessageHandlerFunc) HandleMessage(msg *Message) {
	f(msg)
}

// SendMessage uses the clients transport to encode and send an OSC Message
func (c *Client) SendMessage(msg *Message) error {
	return c.transport.Send(msg)
}

// SendBundle uses the clients transport to encode and send an OSC Bundle
func (c *Client) SendBundle(bun *Bundle) error {
	return c.transport.Send(bun)
}

// SendAndReceiveMessage sends the OSC Message using the clients transport and
// then waits for the response.
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
// it.
func (c *Client) CallMessage(address string, varArg ...any) (*Message, error) {
	return c.SendAndReceiveMessage(&Message{
		Address:   address,
		Arguments: varArg,
	})
}

func (c *Client) listen() {
	var pkg Package
	var err error
	for pkg, err = c.transport.Receive(); err == nil; pkg, err = c.transport.Receive() {
		if pkg.GetType() == PackageTypeMessage {
			m := pkg.(*Message)
			if chi, ok := c.pendingMessageRequests.LoadAndDelete(m.Address); ok {
				ch := chi.(chan *Message)
				ch <- m
				close(ch)
			} else if h, ok := c.messageHandlers[m.Address]; ok {
				h.HandleMessage(m)
			}
		} else if pkg.GetType() == PackageTypeBundle {
			if c.bundleHandler != nil {
				b := pkg.(*Bundle)
				c.bundleHandler.HandleBundle(b)
			}
		}
	}
}