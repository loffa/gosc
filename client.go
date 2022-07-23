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
	messageReceivers       map[string]MessageReceiver
	pendingMessageRequests sync.Map
	bundleReceiver         BundleReceiver
}

type BundleReceiver interface {
	ReceiveBundle(bundle *Bundle)
}

type MessageReceiver interface {
	ReceiveMessage(msg *Message)
}

// MessageReceiverFunc type is an adapter to allow the use of ordinary functions
// as MessageReceiver:s. If f a function with the appropriate signature,
// MessageReceiverFunc(f) is a MessageReceiver that calls f.
type MessageReceiverFunc func(msg *Message)

// ReceiveMessage calls m(msg)
func (m MessageReceiverFunc) ReceiveMessage(msg *Message) {
	m(msg)
}

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
		remote:           remote,
		transport:        trans,
		messageReceivers: map[string]MessageReceiver{},
	}
	go cli.listen()

	return cli, nil
}

// ReceiveMessage adds a MessageHandler for messages on a specific address using
// regexp matching.
func (c *Client) ReceiveMessage(addressPattern string, receiver MessageReceiver) error {
	_, err := regexp.Compile(addressPattern)
	if err != nil {
		return fmt.Errorf("addressPattern is not a valid regexp string: %v", err)
	}
	c.messageReceivers[addressPattern] = receiver
	return nil
}

// ReceiveMessageFunc adds a MessageHandlerFunc for messages on a specific
// address using regexp matching.
func (c *Client) ReceiveMessageFunc(addressPattern string, receiverFunc MessageReceiverFunc) error {
	return c.ReceiveMessage(addressPattern, receiverFunc)
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
// then waits (blocking) for the response to arrive to the listener.
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
				for pattern, h := range c.messageReceivers {
					if c.addressMatches(pattern, m.Address) {
						h.ReceiveMessage(m)
						break
					}
				}
			}
		} else if pkg.GetType() == PackageTypeBundle {
			if c.bundleReceiver != nil {
				b := pkg.(*Bundle)
				c.bundleReceiver.ReceiveBundle(b)
			}
		}
	}
}

func (c *Client) addressMatches(pattern, address string) bool {
	matches, _ := regexp.MatchString(pattern, address)
	return matches
}
