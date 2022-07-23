package gosc

import (
	"testing"
)

func TestHandlerFunc_HandlePackage(t *testing.T) {
	// TODO: Implement me
}

func TestMessageHandlerFunc_HandleMessage(t *testing.T) {
	calledFunction := false
	mh := func(_ *ResponseWriter, msg *Message) {
		calledFunction = true
	}
	MessageHandlerFunc(mh).HandleMessage(nil, &Message{
		Address:   "/test",
		Arguments: []any{},
	})
	if !calledFunction {
		t.Error("Expected MessageHandlerFunc to execute message handler")
	}
}

func TestNewServer(t *testing.T) {
	// TODO: Implement me
}

func TestServer_ListenAndServe(t *testing.T) {
	// TODO: Implement me
}

func TestServer_Shutdown(t *testing.T) {
	// TODO: Implement me
}

func TestServer_listen(t *testing.T) {
	// TODO: Implement me
}
