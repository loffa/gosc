package main

import (
	"github.com/loffa/gosc"
	"log"
)

func main() {
	serv := gosc.NewServer(&gosc.ServerOptions{})
	mux := gosc.NewMux(nil)
	mux.Handle("/hello", HandleHello)

	err := serv.ListenAndServe("127.0.0.1:1234", mux)
	if err != nil {
		log.Println(err)
	}
}

// HandleHello implements the gosc.HandlerFunc interface. It writes a
// response to the client with a string.
func HandleHello(message *gosc.Message, responseWriter *gosc.ResponseWriter) {
	err := responseWriter.Send(&gosc.Message{
		Address:   message.Address,
		Arguments: []any{"World"},
	})
	if err != nil {
		log.Println(err)
	}
}
