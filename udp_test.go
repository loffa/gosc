package gosc

import (
	"testing"
)

func TestNewUDPTransport(t *testing.T) {
	t.Run("correctAddress", func(t *testing.T) {
		_, err := NewUDPTransport("127.0.0.1:1234", 512)
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
	})
	t.Run("malformedAddress", func(t *testing.T) {
		_, err := NewUDPTransport("512.abc:001", 512)
		if err == nil {
			t.Error("expected error but none given")
		}
	})
}

func TestNewUDPListen(t *testing.T) {
	t.Run("correctAddress", func(t *testing.T) {
		_, err := NewUDPListen("127.0.0.1:1234", 512)
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
	})
	t.Run("malformedAddress", func(t *testing.T) {
		_, err := NewUDPListen("512.abc:001", 512)
		if err == nil {
			t.Error("expected error but none given")
		}
	})
}

func Test_transportUDP_Receive(t1 *testing.T) {

}

func Test_transportUDP_Send(t1 *testing.T) {

}
