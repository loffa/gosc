package main

import (
	"fmt"
	"github.com/loffa/gosc"
	"log"
)

// PrefixMiddleware is a middleware that just adds a prefix to messages address before passing it further
// down the handler chain
func PrefixMiddleware(prefix string, next gosc.Handler) gosc.Handler {
	return gosc.HandlerFunc(func(responseWriter *gosc.ResponseWriter, pkg gosc.Package) {
		if msg, ok := pkg.(*gosc.Message); ok {
			msg.Address = fmt.Sprintf("/%s/%s", prefix, msg.Address)
		}
		next.ServePackage(responseWriter, pkg)
	})
}

// MyMux Custom mux that responds with "404" if no message handler was found for given address
// or no bundle handler.
type MyMux struct {
	bundleHandler   gosc.BundleHandler
	messageHandlers map[string]gosc.MessageHandler
}

// ServePackage implements the gosc.Handler interface
func (mux *MyMux) ServePackage(w *gosc.ResponseWriter, pkg gosc.Package) {
	switch x := pkg.(type) {
	case *gosc.Message:
		if handler, ok := mux.messageHandlers[x.Address]; ok {
			handler.HandleMessage(w, x)
			return
		}
		_ = w.Send(&gosc.Message{
			Address:   "/404",
			Arguments: []any{"Not found"},
		})
	case *gosc.Bundle:
		if mux.bundleHandler != nil {
			mux.bundleHandler.HandleBundle(w, x)
			return
		}
		_ = w.Send(&gosc.Message{
			Address:   "/404",
			Arguments: []any{"Not found"},
		})
	}
}

func main() {
	serv := gosc.NewServer(&gosc.ServerOptions{})
	mux := gosc.NewMux(nil)
	mux.HandleMessageFunc("/hello", HandleHello)

	err := serv.ListenAndServe("127.0.0.1:1234", mux)
	if err != nil {
		log.Println(err)
	}
}

// HandleHello implements the gosc.MessageHandler interface. It writes a
// response to the client with a string.
func HandleHello(responseWriter *gosc.ResponseWriter, message *gosc.Message) {
	err := responseWriter.Send(&gosc.Message{
		Address:   message.Address,
		Arguments: []any{"World"},
	})
	if err != nil {
		log.Println(err)
	}
}
