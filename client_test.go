package gosc

import (
	"testing"
)

type testMessageHandler struct {
	handled bool
}

func (t *testMessageHandler) HandleMessage(_ *ResponseWriter, _ *Message) {
	t.handled = true
}

type testBundleHandler struct {
	handled bool
}

func (t *testBundleHandler) HandleBundle(_ *ResponseWriter, _ *Bundle) {
	t.handled = true
}

func TestClient_CallMessage(t *testing.T) {

}

func TestClient_EmitMessage(t *testing.T) {

}

func TestClient_HandleMessage(t *testing.T) {
	cli, _ := NewClient("127.0.0.1:1234")
	t.Run("validPattern", func(t *testing.T) {
		err := cli.HandleMessage("/test", &testMessageHandler{})
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		_, ok := cli.messageHandlers["/test"]
		if !ok {
			t.Error("expected message handler to exist in client")
		}
	})
	t.Run("notValidPattern", func(t *testing.T) {
		err := cli.HandleMessage("[0-9", &testMessageHandler{})
		if err == nil {
			t.Error("expected error but none given")
		}
	})
}

func TestClient_HandleMessageFunc(t *testing.T) {
	cli, _ := NewClient("127.0.0.1:1234")
	err := cli.HandleMessageFunc("/test", func(_ *ResponseWriter, msg *Message) {})
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestClient_SendAndReceiveMessage(t *testing.T) {

}

func TestClient_SendBundle(t *testing.T) {

}

func TestClient_SendMessage(t *testing.T) {

}

func TestClient_addressMatches(t *testing.T) {
	cli, _ := NewClient("127.0.0.1:1234")
	t.Run("match", func(t *testing.T) {
		cli.addressMatches("*", "/test")
	})
	t.Run("notMatch", func(t *testing.T) {
		cli.addressMatches("[0-9]", "/test")
	})
}

func TestClient_listen(t *testing.T) {

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

func TestNewClient(t *testing.T) {
	t.Run("correctClient", func(t *testing.T) {
		cli, err := NewClient("127.0.0.1:1234")
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		if cli.messageHandlers == nil {
			t.Errorf("expected client message messageHandlers to be initialized.")
		}
	})
	t.Run("wrongAddress", func(t *testing.T) {
		_, err := NewClient("127.abc.0.1:1234")
		if err == nil {
			t.Error("expected error but none given")
		}
	})
}
