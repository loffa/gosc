package gosc

import "net"

// A Server defines needed options for running an OSC Server.
type Server struct {
	opts           *ServerOptions
	transport      Transport
	messageHandler Mux
	exiting        chan bool
}

// ServerOptions is the configuration parameters used to create a Server.
type ServerOptions struct {
	// Size for read buffer, defaults to 512
	BufferSize int
}

// NewServer returns a Server with options applied. Options not set get their
// default values.
func NewServer(opts *ServerOptions) *Server {
	if opts.BufferSize == 0 {
		opts.BufferSize = 512
	}

	return &Server{
		opts:    opts,
		exiting: make(chan bool),
	}
}

// ListenAndServe listens on the UDP address specified and then calls handlers
// for incoming packages using the Mux.
//
// ListenAndServe returns error if the address is malformed or can't be opened.
func (s *Server) ListenAndServe(addr string, handler Mux) error {
	trans, err := NewUDPListen(addr, s.opts.BufferSize)
	if err != nil {
		return err
	}
	s.transport = trans
	s.messageHandler = handler
	s.listen()
	return nil
}

// Shutdown waits for all ongoing requests and then shuts down the server and
// releases the port.
func (s *Server) Shutdown() error {
	// TODO: Handle waiting for shutdown to complete
	s.exiting <- true
	return nil
}

func (s *Server) listen() {
	var pkg Package
	var err error
	var src net.Addr
	for pkg, src, err = s.transport.Receive(); err == nil; pkg, src, err = s.transport.Receive() {
		if pkg.GetType() == PackageTypeMessage {
			m := pkg.(*Message)
			s.messageHandler.HandleMessage(src, s.transport, m)
		} else if pkg.GetType() == PackageTypeBundle {
			b := pkg.(*Bundle)
			s.messageHandler.HandleBundle(src, b)
		}
	}
}
