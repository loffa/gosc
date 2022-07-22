package gosc

import (
	"testing"
)

type testMessageHandler struct {
	handled bool
}

func (h *testMessageHandler) HandleMessage(_ *ResponseWriter, _ *Message) {
	h.handled = true
}

type testBundleHandler struct {
	handled bool
}

func (h *testBundleHandler) HandleBundle(_ *ResponseWriter, _ *Bundle) {
	h.handled = true
}

func TestNewMux(t *testing.T) {
	mux := NewMux(nil)
	if mux.messageHandlers == nil {
		t.Error("expected Mux messageHandlers to be initialized")
	}
}

func TestMux_HandlePackage(t *testing.T) {
	t.Run("bundleWithHandler", func(t *testing.T) {
		hand := &testBundleHandler{}
		mux := NewMux(hand)
		mux.HandlePackage(nil, &Bundle{})
		if hand.handled != true {
			t.Error("expected Mux with bundleHandler to handle the bundle")
		}
	})
	t.Run("bundleWithoutHandler", func(t *testing.T) {
		mux := NewMux(nil)
		mux.HandlePackage(nil, &Bundle{})
	})
	t.Run("message", func(t *testing.T) {
		handled := false
		mux := NewMux(nil)
		mux.HandleMessageFunc("/test", func(w *ResponseWriter, m *Message) {
			handled = true
		})
		mux.HandlePackage(nil, &Message{
			Address:   "/test",
			Arguments: []any{},
		})
		if !handled {
			t.Error("expected Mux messageHandler for message to run")
		}
	})
}

func TestDefaultMux_HandleMessage(t *testing.T) {
	mux := NewMux(nil)
	han := &testMessageHandler{}
	mux.HandleMessage("/test", han)
	if len(mux.messageHandlers) != 1 {
		t.Errorf("expected Mux messageHandlers to have 1 entry, got: %d", len(mux.messageHandlers))
	}
}

func TestMux_HandleMessageFunc(t *testing.T) {
	mux := NewMux(nil)
	mux.HandleMessageFunc("/test", func(w *ResponseWriter, msg *Message) {})
	if len(mux.messageHandlers) != 1 {
		t.Errorf("expected Mux messageHandlers to have 1 entry, got: %d", len(mux.messageHandlers))
	}
}

func TestMux_Send(t *testing.T) {
	// TODO: Implement me
}
