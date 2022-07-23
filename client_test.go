package gosc

import (
	"testing"
)

type testMessageReceiver struct {
	received bool
}

func (t *testMessageReceiver) ReceiveMessage(_ *Message) {
	t.received = true
}

func TestClient_CallMessage(t *testing.T) {
	// TODO: Implement me
}

func TestClient_EmitMessage(t *testing.T) {
	// TODO: Implement me
}

func TestClient_ReceiveMessage(t *testing.T) {
	cli, _ := NewClient("127.0.0.1:1234")
	t.Run("validPattern", func(t *testing.T) {
		err := cli.ReceiveMessage("/test", &testMessageReceiver{})
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		_, ok := cli.messageReceivers["/test"]
		if !ok {
			t.Error("expected message handler to exist in client")
		}
	})
	t.Run("notValidPattern", func(t *testing.T) {
		err := cli.ReceiveMessage("[0-9", &testMessageReceiver{})
		if err == nil {
			t.Error("expected error but none given")
		}
	})
}

func TestClient_ReceiveMessageFunc(t *testing.T) {
	cli, _ := NewClient("127.0.0.1:1234")
	err := cli.ReceiveMessageFunc("/test", func(msg *Message) {})
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
}

func TestClient_SendAndReceiveMessage(t *testing.T) {
	// TODO: Implement me
}

func TestClient_SendBundle(t *testing.T) {
	// TODO: Implement me
}

func TestClient_SendMessage(t *testing.T) {
	// TODO: Implement me
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
	// TODO: Implement me
}

func TestNewClient(t *testing.T) {
	t.Run("correctClient", func(t *testing.T) {
		cli, err := NewClient("127.0.0.1:1234")
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
		if cli.messageReceivers == nil {
			t.Errorf("expected client message messageReceivers to be initialized.")
		}
	})
	t.Run("wrongAddress", func(t *testing.T) {
		_, err := NewClient("127.abc.0.1:1234")
		if err == nil {
			t.Error("expected error but none given")
		}
	})
}

func TestMessageReceiverFunc_ReceiveMessage(t *testing.T) {
	// TODO: Implement me
}
