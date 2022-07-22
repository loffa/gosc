package gosc

import "net"

// Mux is an interface used to multiplex requests to the server.
type Mux interface {
	HandleMessage(src net.Addr, t Transport, message *Message)
	HandleBundle(src net.Addr, bundle *Bundle)
}

// DefaultMux is the default multiplexer for messages. Returned by NewMux.
type DefaultMux struct {
	bundleHandler BundleHandler
	handlers      map[string]HandlerFunc
}

// HandlerFunc type is an adapter to allow the use of ordinary functions as
// Handlers.
type HandlerFunc func(message *Message, responseWriter *ResponseWriter)

// NewMux returns the DefaultMux. The bundleHandler can be nil if not handling
// bundles.
func NewMux(bundleHandler BundleHandler) *DefaultMux {
	return &DefaultMux{
		bundleHandler: bundleHandler,
		handlers:      make(map[string]HandlerFunc),
	}
}

// HandleMessage satisfies the Mux interface for DefaultMux
func (m *DefaultMux) HandleMessage(src net.Addr, t Transport, message *Message) {
	handler, ok := m.handlers[message.Address]
	if ok {
		handler(message, &ResponseWriter{
			src:   src,
			trans: t,
		})
	}
}

// HandleBundle satisfies the Mux interface for DefaultMux
func (m *DefaultMux) HandleBundle(_ net.Addr, bundle *Bundle) {
	if m.bundleHandler != nil {
		m.bundleHandler.HandleBundle(bundle)
	}
}

// Handle adds a HandlerFunc to the Mux so that incoming messages on the address
// is dispatched.
func (m *DefaultMux) Handle(address string, handler HandlerFunc) {
	m.handlers[address] = handler
}

// ResponseWriter is used to send responses back to the requesting client on the
// incoming connection.
type ResponseWriter struct {
	src   net.Addr
	trans Transport
}

// Send sends a Package to the client as a response using the Transport of the
// server and the incoming connection.
func (w *ResponseWriter) Send(pkg Package) error {
	return w.trans.Send(pkg, w.src)
}
