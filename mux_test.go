package gosc

import (
	"testing"
)

func TestNewMux(t *testing.T) {
	mux := NewMux(nil)
	if mux.messageHandlers == nil {
		t.Error("expected mux messageHandlers to be initialized")
	}
}

func TestDefaultMux_HandleBundle(t *testing.T) {
	t.Run("withHandler", func(t *testing.T) {
		hand := &testBundleHandler{}
		mux := NewMux(hand)
		mux.ServePackage(nil, &Bundle{})
		if hand.handled != true {
			t.Error("expected default mux with handler to handle the bundle")
		}
	})
	t.Run("withoutHandler", func(t *testing.T) {
		mux := NewMux(nil)
		mux.ServePackage(nil, &Bundle{})
	})
}

func TestDefaultMux_HandleMessage(t *testing.T) {
	handled := false
	mux := NewMux(nil)
	mux.HandleMessageFunc("/test", func(w *ResponseWriter, m *Message) {
		handled = true
	})
	mux.ServePackage(nil, &Message{
		Address:   "/test",
		Arguments: []any{},
	})
	if !handled {
		t.Error("expected default mux handler for message to run")
	}
}

func TestDefaultMux_Handle(t *testing.T) {
	mux := NewMux(nil)
	mux.HandleMessageFunc("/test", func(w *ResponseWriter, msg *Message) {})
	if len(mux.messageHandlers) != 1 {
		t.Errorf("expected default mux messageHandlers to have 1 entry, got: %d", len(mux.messageHandlers))
	}
}
