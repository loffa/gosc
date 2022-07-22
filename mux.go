package gosc

import "net"

// Mux is the default multiplex handler for messages and bundles. Returned by NewMux.
type Mux struct {
	bundleHandler   BundleHandler
	messageHandlers map[string]MessageHandler
}

func (m *Mux) HandlePackage(writer *ResponseWriter, pkg Package) {
	switch x := pkg.(type) {
	case *Message:
		if handler, ok := m.messageHandlers[x.Address]; ok {
			handler.HandleMessage(writer, x)
		}
	case *Bundle:
		if m.bundleHandler != nil {
			m.bundleHandler.HandleBundle(writer, x)
		}
	}
}

func (m *Mux) HandleMessage(addr string, handler MessageHandler) {
	m.messageHandlers[addr] = handler
}

func (m *Mux) HandleMessageFunc(addr string, handlerFunc MessageHandlerFunc) {
	m.messageHandlers[addr] = handlerFunc
}

// NewMux returns the Mux. The bundleHandler can be nil if not handling
// bundles.
func NewMux(bundleHandler BundleHandler) *Mux {
	return &Mux{
		bundleHandler:   bundleHandler,
		messageHandlers: make(map[string]MessageHandler),
	}
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
