package gosc

import (
	"testing"
)

func TestNewMux(t *testing.T) {
	mux := NewMux(nil)
	if mux.handlers == nil {
		t.Error("expected mux handlers to be initialized")
	}
}

func TestDefaultMux_HandleBundle(t *testing.T) {
	t.Run("withHandler", func(t *testing.T) {
		hand := &testBundleHandler{}
		mux := NewMux(hand)
		mux.HandleBundle(nil, &Bundle{})
		if hand.handled != true {
			t.Error("expected default mux with handler to handle the bundle")
		}
	})
	t.Run("withoutHandler", func(t *testing.T) {
		mux := NewMux(nil)
		mux.HandleBundle(nil, &Bundle{})
	})
}

func TestDefaultMux_HandleMessage(t *testing.T) {
	handled := false
	mux := NewMux(nil)
	mux.Handle("/test", func(m *Message, w *ResponseWriter) {
		handled = true
	})
	mux.HandleMessage(nil, nil, &Message{
		Address:   "/test",
		Arguments: []any{},
	})
	if !handled {
		t.Error("expected default mux handler for message to run")
	}
}

func TestDefaultMux_Handle(t *testing.T) {
	mux := NewMux(nil)
	mux.Handle("/test", func(m *Message, w *ResponseWriter) {})
	if len(mux.handlers) != 1 {
		t.Errorf("expected default mux handlers to have 1 entry, got: %d", len(mux.handlers))
	}
}
