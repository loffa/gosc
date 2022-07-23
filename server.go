package gosc

import "net"

// PackageHandler provides an interface dealing with any OSC package
type PackageHandler interface {
	HandlePackage(writer *ResponseWriter, pkg Package)
}

// A MessageHandler is called when messages on the specified address is received.
type MessageHandler interface {
	HandleMessage(writer *ResponseWriter, msg *Message)
}

// A BundleHandler is called when a Bundle is received.
type BundleHandler interface {
	HandleBundle(writer *ResponseWriter, bundle *Bundle)
}

// HandlerFunc type is an adapter to allow the use of ordinary functions as
// Handlers.
type HandlerFunc func(responseWriter *ResponseWriter, pkg Package)

// The MessageHandlerFunc type is an adapter to allow the use of ordinary
// functions as MessageHandler:s. If f is a function with the appropriate
// signature, MessageHandlerFunc(w, f) is a MessageHandler that calls f.
type MessageHandlerFunc func(writer *ResponseWriter, msg *Message)

// HandleMessage calls m(w, msg)
func (m MessageHandlerFunc) HandleMessage(w *ResponseWriter, msg *Message) {
	m(w, msg)
}

// HandlePackage is the HandlerFunc:s implementation of the PackageHandler
// interface.
func (h HandlerFunc) HandlePackage(writer *ResponseWriter, pkg Package) {
	h(writer, pkg)
}

// A Server defines needed options for running an OSC Server. A Server is
// created using the NewServer method.
type Server struct {
	opts           *ServerOptions
	transport      Transport
	packageHandler PackageHandler
	exiting        chan bool
}

// ServerOptions is the configuration parameters used to create a Server.
type ServerOptions struct {
	// Size for read buffer, defaults to 512
	BufferSize int
}

// NewServer initializes and returns a Server with options applied. Options not
// set get their default values.
func NewServer(opts *ServerOptions) *Server {
	if opts.BufferSize == 0 {
		opts.BufferSize = 512
	}

	return &Server{
		opts:    opts,
		exiting: make(chan bool),
	}
}

// ListenAndServe listens on the UDP address specified and then calls
// the PackageHandler for incoming packages.
//
// ListenAndServe returns error if the address is malformed or can't be opened.
func (s *Server) ListenAndServe(addr string, handler PackageHandler) error {
	trans, err := NewUDPListen(addr, s.opts.BufferSize)
	if err != nil {
		return err
	}
	s.transport = trans
	s.packageHandler = handler
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
		rw := &ResponseWriter{
			src:   src,
			trans: s.transport,
		}
		s.packageHandler.HandlePackage(rw, pkg)
	}
}
